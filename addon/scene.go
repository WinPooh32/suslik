package addon

import "github.com/WinPooh32/suslik"

// Test that Object implements Generic interface.
var _ = func(g Generic) struct{} { return struct{}{} }(&Object{})

type Generic interface {
	Update(parent *Object, dt float32)
	Render(parent *Object, batch *suslik.Batch)
	Close(parent *Object)

	UpdateChildren(parent *Object, dt float32)
	RenderChildren(parent *Object, batch *suslik.Batch)
	CloseChildren(parent *Object)

	Base() *Object
}

type Object struct {
	ID       uint
	ParentID uint
	Pos      suslik.Point
	Children []Generic
	Deleted  bool
}

func (object *Object) Update(parent *Object, dt float32) {}

func (object *Object) Render(parent *Object, batch *suslik.Batch) {}

func (object *Object) Close(parent *Object) {}

func (object *Object) UpdateChildren(parent *Object, dt float32) {
	var cleanup []int

	for i, generic := range object.Children {
		child := generic.Base()
		if !child.Deleted {
			generic.Update(object, dt)
			generic.UpdateChildren(object, dt)
		} else {
			cleanup = append(cleanup, i)
		}
	}

	if len(cleanup) > 0 {
		var last = len(object.Children) - 1

		if last > 0 {
			for _, i := range cleanup {
				object.Children[i] = object.Children[last]
				last--
			}
		} else {
			object.Children = object.Children[:0]
		}

		object.Children = object.Children[:last+1]
	}
}

func (object *Object) RenderChildren(parent *Object, batch *suslik.Batch) {
	for _, generic := range object.Children {
		child := generic.Base()
		if !child.Deleted {
			generic.Render(object, batch)
			generic.RenderChildren(object, batch)
		}
	}
}

func (object *Object) CloseChildren(parent *Object) {
	for _, generic := range object.Children {
		child := generic.Base()
		if !child.Deleted {
			generic.Close(object)
			generic.CloseChildren(object)
		}
	}
}

func (object *Object) Base() *Object {
	return object
}

type ObjectConstructor func() interface{}

type World struct {
	lastID       uint
	registry     map[uint]interface{}
	constructors map[string]func() interface{}
}

func NewWorld() *World {
	var world = World{
		lastID:       0,
		registry:     map[uint]interface{}{},
		constructors: map[string]func() interface{}{},
	}

	world.RegisterConstructor("__root", func() interface{} { return new(Object) })

	return &world
}

func (world *World) RegisterConstructor(name string, constructor ObjectConstructor) {
	world.constructors[name] = constructor
}

func (world *World) Build(name string, parentID uint) (interface{}, bool) {
	construct, ok := world.constructors[name]
	if !ok {
		return nil, false
	}

	var abstract = construct()
	var object *Object
	var id uint

	if generic, ok := abstract.(Generic); ok && generic != nil {
		world.lastID++
		id = world.lastID
		object = generic.Base()
		object.ID = world.lastID
		object.ParentID = parentID

		if parent, ok := world.registry[parentID]; ok {
			parentBase := parent.(Generic).Base()
			parentBase.Children = append(parentBase.Children, generic)
		} else if parentID > 0 {
			return nil, false
		}
	} else {
		return nil, false
	}

	world.registry[id] = abstract
	return abstract, true
}

func (world *World) Move(srcID, dstID uint) {
	abstract, ok := world.registry[srcID]
	if !ok {
		return
	}
	objectSrc, ok := abstract.(*Object)
	if !ok {
		return
	}

	abstract, ok = world.registry[dstID]
	if !ok {
		return
	}
	parent, ok := abstract.(*Object)
	if !ok {
		return
	}

	abstract, ok = world.registry[dstID]
	if !ok {
		return
	}
	objectDst, ok := abstract.(*Object)
	if !ok {
		return
	}

	var last = len(parent.Children) - 1
	if last > 0 {
		var i int = -1
		for j, child := range parent.Children {
			if child.Base().ID == objectSrc.ID {
				i = j
			}
		}
		if i < 0 {
			return
		}

		parent.Children[i] = parent.Children[last]
		parent.Children = parent.Children[:last+1]

	} else if parent.Children[0].Base().ID == objectSrc.ID {
		parent.Children = parent.Children[:0]

	} else {
		return
	}

	objectSrc.ParentID = objectDst.ID
	objectDst.Children = append(objectDst.Children, objectSrc)
}

func (world *World) Delete(id uint) {
	if abstract, ok := world.registry[id]; ok {
		if generic, ok := abstract.(Generic); ok {
			object := generic.Base()
			object.Deleted = true
			for _, generic := range object.Children {
				child := generic.Base()
				child.Deleted = true
				world.Delete(child.ID)
			}
		}
	}
	delete(world.registry, id)
}

type Scene struct {
	suslik.Game
	world *World
	root  *Object
	batch *suslik.Batch
}

func NewScene(world *World) *Scene {
	var root, _ = world.Build("__root", 0)
	return &Scene{
		world: world,
		root:  root.(*Object),
		batch: nil,
	}
}

func (scene *Scene) World() *World {
	return scene.world
}

func (scene *Scene) Root() uint {
	return scene.root.ID
}

func (scene *Scene) Preload() {

}

func (scene *Scene) Setup() {
	scene.batch = suslik.NewBatch(suslik.Width(), suslik.Height())
}

func (scene *Scene) Update(dt float32) {
	scene.root.UpdateChildren(nil, dt)
}

func (scene *Scene) Render() {
	scene.batch.Begin()
	scene.root.RenderChildren(nil, scene.batch)
	scene.batch.End()
}

func (scene *Scene) Resize(w, h float32) {
	scene.batch.SetProjection(w, h)
}

func (scene *Scene) Close() {
	scene.root.CloseChildren(nil)
	scene.world = nil
}
