package glvox

type Tracer interface {
	Trace(pos, dir Vec3) (dest Vec3, hit bool)
}

type Getter interface {
	Get(x, y, z int) (val, size int)
}

type Setter interface {
	Set(x, y, z int, v int)
}

type GetSetter interface {
	Getter
	Setter
}

type Sized interface {
	Size() Size
}

type SizedGetter interface {
	Sized
	Getter
}

type Size struct {
	W, H, D int
}

type Vox struct {
	Dist Vec3
	Size float32
	Value float32
}

