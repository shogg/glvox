#version 140

in  vec2 texCoord;
out vec4 fragColor;

uniform float time;
uniform isamplerBuffer voxels;
uniform int size;
uniform bool shadowOff;

uniform struct {
	vec3 pos, dir, up, left;
} cam;

vec3 rotateX(vec3 p, float a)
{
    float sa = sin(a);
    float ca = cos(a);
    vec3 r;
    r.x = p.x;
    r.y = ca*p.y - sa*p.z;
    r.z = sa*p.y + ca*p.z;
    return r;
}

vec3 rotateY(vec3 p, float a)
{
    float sa = sin(a);
    float ca = cos(a);
    vec3 r;
    r.x = ca*p.x + sa*p.z;
    r.y = p.y;
    r.z = -sa*p.x + ca*p.z;
    return r;
}

vec3 face(vec3 p, vec3 d, float off, float size, int axis, out bool hit)
{
	float f;
	switch(axis) {
	case 0:
		if(d.x == 0.0) { hit = false; return p; }
		f = (off - p.x) / d.x;
		break;
	case 1:
		if(d.y == 0.0) { hit = false; return p; }
		f = (off - p.y) / d.y;
		break;
	case 2:
		if(d.z == 0.0) { hit = false; return p; }
		f = (off - p.z) / d.z;
		break;
	}

	vec3 pos = p + f*d;

	switch(axis) {
	case 0:
		hit = pos.y >= 0.0 && pos.y <= size &&
			pos.z >= 0.0 && pos.z <= size;
		pos.x = int(pos.x);
		break;
	case 1:
		hit = pos.x >= 0.0 && pos.x <= size &&
			pos.z >= 0.0 && pos.z <= size;
		pos.y = int(pos.y);
		break;
	case 2:
		hit = pos.x >= 0.0 && pos.x <= size &&
			pos.y >= 0.0 && pos.y <= size; 
		pos.z = int(pos.z);
		break;
	}

	return pos;
}

vec3 boundingBox(vec3 p, vec3 d, out bool hit)
{
	float size = ::size;

	// inside the box
	if(	p.x >= 0.0 && p.x <= size &&
		p.y >= 0.0 && p.y <= size &&
		p.z >= 0.0 && p.z <= size) {

		hit = true;
		return p;
	}

	vec3 pos;

	if(d.x > 0.0) {
		pos = face(p, d, 0.0, size, 0, hit);
		if(hit) { return pos; }
	} else {
		pos = face(p, d, size, size, 0, hit);
		if(hit) { return pos; }
	}

	if(d.y > 0.0) {
		pos = face(p, d, 0.0, size, 1, hit);
		if(hit) { return pos; }
	} else {
		pos = face(p, d, size, size, 1, hit);
		if(hit) { return pos; }
	}

	if(d.z > 0.0) {
		pos = face(p, d, 0.0, size, 2, hit);
		if(hit) { return pos; }
	} else {
		pos = face(p, d, size, size, 2, hit);
		if(hit) { return pos; }
	}

	return p;
}

struct vox {
	vec3 dist;
	float size;
	float alpha;
	int x, y, z;
	int steps;
};

int get(int x, int y, int z, out int size)
{
	const int max = 8;
	const int max2 = 8;

	size = 1;
	if(	x < 0 || x >= max2 ||
		y < 0 || y >= max2 ||
		z < 0 || z >= max2) { return 0; }

	if(+x+y+z < max) { return 1; }

	return 0;
}

int octree(int x, int y, int z, out int size)
{
	size = ::size;

	if(x < 0 || x >= size || y < 0 || y >= size || z < 0 || z >= size) {
		return 0;
	}

	int steps = 0;
	int i = 0, off = 0;
	while(size > 1) {
		steps++;

		size >>= 1;
		off = 0;

		if(z >= size) { off += 4; z -= size; }
		if(y >= size) { off += 2; y -= size; }
		if(x >= size) { off += 1; x -= size; }

		i = texelFetch(voxels, (i<<3) + off).r;
		if(i <= 0) {
			return -i;
		}
	}

	return 0;
}

vox voxel(vec3 pos, vec3 dir)
{
	int x = int(pos.x);
	int y = int(pos.y);
	int z = int(pos.z);

	if(dir.x < 0.0 && abs(fract(pos.x)) == 0.0) { x--; }
	if(dir.y < 0.0 && abs(fract(pos.y)) == 0.0) { y--; }
	if(dir.z < 0.0 && abs(fract(pos.z)) == 0.0) { z--; }

	int s;
	int val = octree(x, y, z, s);

	float size = float(s) * .5;
	vec3 center = vec3(x/s*s, y/s*s, z/s*s) + vec3(size);
	vec3 dist = pos - center;
	vox v = vox(dist, size, float(val), x/s*s, y/s*s, z/s*s, val);

	return v;
}

vec3 trace(vec3 o, vec3 d, out vec3 n, out bool hit, out vox v)
{
	const vec3 nx = vec3(1.0, .0, .0);
	const vec3 ny = vec3(.0, 1.0, .0);
	const vec3 nz = vec3(.0, .0, 1.0);

	vec3 s = vec3(1.0);
	if(d.x < 0.0) { s.x = -1.0; }
	if(d.y < 0.0) { s.y = -1.0; }
	if(d.z < 0.0) { s.z = -1.0; }
/*
	vec3 pos = boundingBox(o, d, hit);
	if(!hit) { return o; }
*/
	vec3 pos = o;
	hit = false;

	const int maxSteps = 51;
	for(int i = 0; i < maxSteps; i++) {
		v = voxel(pos, s);
		v.steps = i;
		if(v.alpha > 0.0) {
			n *= -s;
			hit = true;
			return pos;
		}

		vec3 f = s*vec3(v.size) - v.dist;
		f.x /= d.x; f.y /= d.y; f.z /= d.z;

		float fmin = 100.0;
		if(f.x > 0.0 && f.x < fmin) { fmin = f.x; n = nx; }
		if(f.y > 0.0 && f.y < fmin) { fmin = f.y; n = ny; }
		if(f.z > 0.0 && f.z < fmin) { fmin = f.z; n = nz; }

		pos = pos + d*fmin;
	}

	return o;
}

vec3 background(vec3 d)
{
	return mix(vec3(1.0), vec3(0.0, 0.25, 1.0), d.y) - .4;
}
/*
float ambientOcclusion(vec3 pos, vox v)
{

}
*/
float shadow(vec3 pos, vec3 lightPos)
{
	vec3 ray = lightPos - pos;
	ray = normalize(ray);

	vec3 n; bool hit; vox v;
	vec3 p = trace(pos, ray, n, hit, v);
	if(hit) { return 0.3; }
	return 1.0;
}

vec3 shade(vec3 pos, vec3 n, vec3 eyePos, vox voxel)
{
	//n = normalize(voxel.dist);

	const vec3 color = vec3(2.0, 1.5, 1.2);
    const vec3 lightPos = vec3( 0000.0,-0000.0, 0000.0);
    const float shininess = 130.0;

    vec3 l = normalize(lightPos - pos);
    vec3 v = normalize(eyePos - pos);
    vec3 h = normalize(v + l);
    float diff = dot(n, l);
    float spec = max(0.0, pow(dot(n, h), shininess)) * float(diff > 0.0);
    diff = 0.6+0.4*diff;

	float shadow = shadowOff ? 1.0 : shadow(pos, lightPos);

	//if(voxel.size > 1) { return vec3(0, voxel.size*0.001, 0); }
/*
	vec3 green = vec3(0.0, 1.6, 0.0);
	vec3 red = vec3(1.0, 0.0, 0.0);
	vec3 blue = vec3(0.0, 0.0, 0.6);
	vec3 yellow = vec3(1.0, 1.0, 0.0);
	vec3 rgb = mix(green, blue, voxel.steps*0.045);
	if(voxel.steps > 20) {
		rgb = mix(red, yellow, (voxel.steps-20)*0.09);
	}
*/
	//vec3 rgb = shadow * abs(voxel.dist)*1.7;
	//vec3 rgb = vec3(0, voxel.alpha*1.5, 0)*1.7;
	//if(voxel.alpha == 8) { discard; }
	//vec3 rgb = vec3(voxel.size*2.0);
	//if(length(voxel.dist) > 0.55) { rgb = vec3(0, 0, 0.0); }
	vec3 rgb = shadow * vec3(voxel.x, voxel.y, voxel.z)*.0008;
	return rgb * diff; // + spec;
}

void main(void)
{
	vec3 o = cam.pos + vec3(-1260.0, 1000.0, 500.0);
	vec3 d = normalize(
		cam.dir*2.0 + cam.left*texCoord.x + cam.up*texCoord.y);

#if 1
    float a = sin(time*0.4)*0.001+1.8099;
    o = rotateY(o, a);
    d = rotateY(d, a);
#endif

	bool hit = false;
	vec3 n; vox v;
	vec3 p = trace(o, d, n, hit, v);

	vec3 rgb;
	if(hit) {
		rgb = shade(p, n, o, v);
	} else {
		rgb = background(d);
	}

	fragColor = vec4(rgb, 1.0);
}

