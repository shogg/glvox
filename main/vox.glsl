//     +---+---+
//    /   /   /|
//   +---+---+ +
//  /   /   /|/|
// +---+---+ + +
// |   |   |/|/
// +---+---+ +
// |   |   |/
// +---+---+

varying vec2 texCoord;
uniform float time;
uniform sampler1D voxels;

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
		break;
	case 1:
		hit = pos.x >= 0.0 && pos.x <= size &&
			pos.z >= 0.0 && pos.z <= size;
		break;
	case 2:
		hit = pos.x >= 0.0 && pos.x <= size &&
			pos.y >= 0.0 && pos.y <= size; 
		break;
	}

	return pos;
}

vec3 boundingBox(vec3 p, vec3 d, out bool hit)
{
	const float size = 2.0;

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

vec3 background(vec3 d)
{
	if(d.y > 0.0) {
		return mix(vec3(1.0), vec3(0.0, 0.25, 1.0), d.y);
	}
	return mix(vec3(3.5), vec3(0.1, .6, 0.3), -d.y*2.5);
}

vec3 shade(vec3 p, vec3 n)
{
	vec3 rgb = vec3(p.x*.7, p.y*.4, p.z*.2);
	return rgb;
}

void main(void)
{
	vec3 o = vec3(0.0, .5, -5.5);
	vec3 d = normalize(vec3(texCoord.x, texCoord.y, 2.0));

#if 1
    float a = sin(time*0.4)*3.5;
    o = rotateY(o, a);
    d = rotateY(d, a);
#endif

	bool hit = false;
	vec3 p = boundingBox(o, d, hit);

	vec3 rgb;
	if(hit) {
		vec3 n = vec3(1.0);
		rgb = shade(p, n);
	} else {
		rgb = background(d);
	}

	gl_FragColor = vec4(rgb, 1.0);
}

