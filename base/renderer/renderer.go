package renderer

import (
	"gotracer/base/scene"
	"gotracer/base/scene/objects"
	"gotracer/base/scene/primitives/triangle"
	"gotracer/base/vector"
	"image"
	"image/color"
	"log"
	"math"
)

type Ray struct {
	vector.Vec4f
	Pixel vector.Vec2i
	T     float32
}

func triangleIntersect(camera objects.Camera, ray *Ray, tri triangle.Triangle) (intersectionPoint vector.Vec4f, hit bool) {
	e1 := tri.Vertices[1].Subtract(&(tri.Vertices[0].Vec4f))
	e2 := tri.Vertices[2].Subtract(&(tri.Vertices[0].Vec4f))

	N := e1.Cross(e2)
	cameraRay := camera.Position.Add(&(ray.Vec4f))
	dotRayTri := N.Dot(cameraRay)
	if dotRayTri == 0 || dotRayTri > 0 {
		hit = false
		return
	}

	h := ray.Cross(e2)
	a := e1.Dot(h)

	if a > -0.00001 && a < 0.00001 {
		hit = false
		return
	}

	f := 1 / a
	s := camera.Position.Subtract(&(tri.Vertices[0].Vec4f))
	u := f * s.Dot(h)

	if u < 0.0 || u > 1.0 {
		hit = false
		return
	}

	q := s.Cross(e1)
	v := f * ray.Dot(q)

	if v < 0.0 || u+v > 1.0 {
		hit = false
		return
	}

	t := f * e2.Dot(q)

	if t > 0.00001 {
		if t < ray.T {
			ray.T = t
		} else {
			//log.Println("Not the closest hit, discarding.")
			hit = false
			return
		}
		hit = true
	} else {
		hit = false
		return
	}

	scaledRay := ray.Scale(t)
	intersectionPoint = camera.Position.Add(&scaledRay)

	return
}

func getPrincipalRays(world *scene.Scene, w, h uint32) (rays []*Ray) {
	log.Println("Creating principal rays")
	for i := 0; uint32(i) < w; i++ {
		for j := 0; uint32(j) < h; j++ {
			vec := new(Ray)
			vec.X = float32(i) - float32(w)/2
			vec.Y = float32(j) - float32(h)/2
			vec.Z = (-1) * world.Camera.Position.Z
			vec.W = 1.0
			vec.Normalize()
			vec.T = math.MaxFloat32
			vec.Pixel.X = int32(i)
			vec.Pixel.Y = int32(j)

			rays = append(rays, vec)
		}
	}
	return
}

func intersect(world *scene.Scene, rays []*Ray, canvas *image.RGBA) {
	for _, ray := range rays {
		//var latest triangle.Triangle
		for _, tri := range world.Triangles {
			_, hit := triangleIntersect(world.Camera, ray, tri)
			if hit {
				//log.Println("Hit! Setting color of pixel", ray.Pixel.X, ray.Pixel.Y)
				canvas.Set(int(ray.Pixel.X), int(ray.Pixel.Y), color.RGBA{0xff, 0, 0, 0xff})
				//latest = tri
				//log.Println("Not the closest hit, discarding.")log.Println("Triangle ", tri, "hit")
			} else {
				//log.Println("Triangle ", tri, "not hit")
			}
		}
	}
	return
}

func Render(world *scene.Scene, w, h uint32, canvas *image.RGBA) {
	for {

		log.Println("Starting tracing rays")
		world.Lock()
		//generate principal ray(s)
		principalRays := getPrincipalRays(world, w, h)
		//intersect ray(s)
		intersect(world, principalRays, canvas)
		//get color
		world.Unlock()
		log.Println("Done tracing rays")
		break
	}
	return
}
