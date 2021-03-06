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
	ivec3 coord;
	int steps;
};

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

	if(dir.x < 0.0 && (fract(pos.x)) == 0.0) { x--; }
	if(dir.y < 0.0 && (fract(pos.y)) == 0.0) { y--; }
	if(dir.z < 0.0 && (fract(pos.z)) == 0.0) { z--; }

	int s;
	int val = octree(x, y, z, s);

	ivec3 coord = ivec3(x, y, z)/s*s;

	float size = float(s) * .5;
	vec3 center = coord + vec3(size);
	vec3 dist = pos - center;
	vox v = vox(dist, size, float(val), coord, val);

	return v;
}

vec3 normal(vec3 pos, vec3 dir)
{
	int x = int(pos.x);
	int y = int(pos.y);
	int z = int(pos.z);

	if(dir.x < 0.0 && (fract(pos.x)) == 0.0) { x--; }
	if(dir.y < 0.0 && (fract(pos.y)) == 0.0) { y--; }
	if(dir.z < 0.0 && (fract(pos.z)) == 0.0) { z--; }

	int s;
	int density = octree(x, y, z, s);
	float dx0 = octree(x-1, y, z, s) - density;
	float dy0 = octree(x, y-1, z, s) - density;
	float dz0 = octree(x, y, z-1, s) - density;
	vec3 n0 = normalize(vec3(dx0, dy0, dz0));

	float dx1 = density - octree(x+1, y, z, s);
	float dy1 = density - octree(x, y+1, z, s);
	float dz1 = density - octree(x, y, z+1, s);
	vec3 n1 = normalize(vec3(dx1, dy1, dz1));

	vec3 n = mix(n0, n1, fract(pos));
	return n;

/*
	float dx = octree(x-1, y, z, s) - octree(x+1, y, z, s);
	float dy = octree(x, y-1, z, s) - octree(x, y+1, z, s);
	float dz = octree(x, y, z-1, s) - octree(x, y, z+1, s);

	vec3 n = normalize(vec3(dx, dy, dz));
	return n;
*/
}

bool box(vec3 pos, vec3 dir, out float t0, out float t1)
{
	int size = ::size;

	vec3 boxMin = vec3(0.0);
	vec3 boxMax = vec3(size);

	vec3 invR = 1.0 / dir;
	vec3 tbot = invR * (boxMin - pos);
	vec3 ttop = invR * (boxMax - pos);
	vec3 tmin = min(ttop, tbot);
	vec3 tmax = max(ttop, tbot);

	vec2 t = max(tmin.xx, tmin.yz);
	t0 = max(t.x, t.y);
	t = min(tmax.xx, tmax.yz);
	t1 = min(t.x, t.y);

	return t0 <= t1;
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

	float t0, t1;
	hit = box(o, d, t0, t1);
	if(!hit) { return o; }

	vec3 pos = o + d * t0;
	vec3 posEnd = o + d * t1;
*/
	vec3 pos = o;
	hit = false;

	const int maxSteps = 31;
	for(int i = 0; i < maxSteps; i++) {
		v = voxel(pos, s);
		v.steps = i;
		if(v.alpha > 0.0) {
			n *= -s;
			//n = normal(pos, s);
			hit = true;
			return pos;
		}

		vec3 f = (s*vec3(v.size) - v.dist) / d;

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

	vec3 rgb = shadow * vec3(voxel.coord.x, voxel.coord.y, voxel.coord.z)*.0008;
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

