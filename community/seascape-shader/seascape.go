package main

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	// Set up window configs
	cfg := pixelgl.WindowConfig{ // Default: 1024 x 768
		Title:  "Golang Seascape from Shadertoy",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	camVector := win.Bounds().Center()

	bounds := win.Bounds()
	bounds.Max = bounds.Max.ScaledXY(pixel.V(1.0, 1.0))

	// I am putting all shader example initializing stuff here for
	// easier reference to those learning to use this functionality
	fragSource, err := LoadFileToString("shaders/seascape.glsl")
	if err != nil {
		panic(err)
	}

	var uMouse mgl32.Vec4
	var uTime float32

	canvas := win.Canvas()
	uResolution := mgl32.Vec2{float32(win.Bounds().W()), float32(win.Bounds().H())}

	EasyBindUniforms(canvas,
		"u_resolution", &uResolution,
		"u_time", &uTime,
		"u_mouse", &uMouse,
	)

	canvas.SetFragmentShader(fragSource)

	start := time.Now()

	// Game Loop
	for !win.Closed() {
		uTime = float32(time.Since(start).Seconds())
		mpos := win.MousePosition()
		uMouse[0] = float32(mpos.X)
		uMouse[1] = float32(mpos.Y)

		win.Clear(colornames.Black)

		// Drawing to the screen
		canvas.Draw(win, pixel.IM.Moved(camVector))

		win.Update()
	}

}

func main() {
	pixelgl.Run(run)
}

var cloudsFragmentShader = `
#version 330 core

#ifdef GL_ES
precision highp float;
#endif

#define HOW_CLOUDY 0.2
#define SHADOW_THRESHOLD 0.4
#define SHADOW 0.3
#define SUBSURFACE 1.0
#define WIND_DIRECTION 0.3
#define TIME_SCALE 0.6
#define SCALE 0.1
//#define ENABLE_SHAFTS
in vec2 texcoords;
out vec4 fragColor;
mat2 RM = mat2(cos(WIND_DIRECTION), -sin(WIND_DIRECTION), sin(WIND_DIRECTION), cos(WIND_DIRECTION));
uniform float u_time;
uniform vec2 u_mouse;
//uniform vec2 u_resolution;
uniform vec4 u_texbounds;
uniform sampler2D u_texture;
uniform vec2 u_gopherpos;

float hash( float n )
{
    return fract(sin(n)*758.5453);
}

float noise( in vec3 x )
{
    vec3 p = floor(x);
    vec3 f = fract(x);
    float n = p.x + p.y*57.0 + p.z*800.0;
    float res = mix(mix(mix( hash(n+  0.0), hash(n+  1.0),f.x), mix( hash(n+ 57.0), hash(n+ 58.0),f.x),f.y),
            mix(mix( hash(n+800.0), hash(n+801.0),f.x), mix( hash(n+857.0), hash(n+858.0),f.x),f.y),f.z);
    return res;
}

float fbm( vec3 p )
{
    float f = 0.0;
    f += 0.50000*noise( p ); p = p*2.02;
    f -= 0.25000*noise( p ); p = p*2.03;
    f += 0.12500*noise( p ); p = p*3.01;
    f += 0.06250*noise( p ); p = p*3.04;
    f += 0.03500*noise( p ); p = p*4.01;
    f += 0.01250*noise( p ); p = p*4.04;
    f -= 0.00125*noise( p );
    return f/0.784375;
}

float cloud(vec3 p)
{
    p-=fbm(vec3(p.x,p.y,0.0)*0.5)*1.25;
    float a = min((fbm(p*3.0)*2.2-1.1), 0.0);
    return a*a;
}

float shadow = 1.0;


float clouds(vec2 p){
    float ic = cloud(vec3(p * 2.0, u_time*0.01 * TIME_SCALE)) / HOW_CLOUDY;
    float init = smoothstep(0.1, 1.0, ic) * 5.0;
    shadow = smoothstep(0.0, SHADOW_THRESHOLD, ic) * SHADOW + (1.0 - SHADOW);
    init = (init * cloud(vec3(p * (6.0), u_time*0.01 * TIME_SCALE)) * ic);
    init = (init * (cloud(vec3(p * (11.0), u_time*0.01 * TIME_SCALE))*0.5 + 0.4) * init);
    return min(1.0, init);
}
//uniform sampler2D bb;
float cloudslowres(vec2 p){
    return 1.0 - (texture(u_texture, p).a - 0.9) * 10.0;
}

vec2 ratio = vec2(1.0, 1.0);

vec4 getresult(){
	vec2 uvmouse = (u_mouse/(texcoords - u_texbounds.xy));
	vec2 t = (texcoords - u_texbounds.xy) / u_texbounds.zw;
    //vec2 surfacePosition = ((( t ) * vec2(u_gopherpos.x , u_gopherpos.y)) * 2.0 - 1.0)*SCALE;
	vec2 surfacePosition = t+u_gopherpos*10.0;
	vec2 position = ( surfacePosition * SCALE);
	vec2 sun = (uvmouse.xy * vec2(texcoords.x / texcoords.y, 1.0)*2.0-1.0) * SCALE;

    float dst = distance(sun * ratio, position * ratio);
    float suni = pow(dst + 1.0, -10.0);
    float shaft =0.0;
    float st = 0.05;
    float w = 1.0;
    vec2 dir = sun - position;
    float c = clouds(position);
    #ifdef ENABLE_SHAFTS
    for(int i=0;i<50;i++){
        float occl = cloudslowres(clamp((t) + dir * st, 0.0, 1.0));
        w *= 0.99;
        st *= 1.05;
        shaft += max(0.0, (1.0 - occl)) * w;
    }
    #endif
    shadow = min(1.0, shadow + suni * suni * 0.2 * SUBSURFACE);
    suni *= (shaft * 0.03);
    return vec4(pow(mix(vec3(shadow), pow(vec3(0.23, 0.33, 0.48), vec3(2.2)) + suni, c), vec3(1.0/2.2)), c*0.1 + 0.9);     
}

void main( void ) {
    fragColor = getresult().rgba;
}
`
