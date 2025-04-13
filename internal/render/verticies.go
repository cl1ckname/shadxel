package render

type Vertex struct {
	Position [3]float32
	Normal   [3]float32
}

// Helper to add a face
func addFace(vertices *[]float32, normal [3]float32, corners [4][3]float32) {
	// Triangle 1
	*vertices = append(*vertices,
		corners[0][0], corners[0][1], corners[0][2], normal[0], normal[1], normal[2],
		corners[1][0], corners[1][1], corners[1][2], normal[0], normal[1], normal[2],
		corners[2][0], corners[2][1], corners[2][2], normal[0], normal[1], normal[2],
	)
	// Triangle 2
	*vertices = append(*vertices,
		corners[2][0], corners[2][1], corners[2][2], normal[0], normal[1], normal[2],
		corners[3][0], corners[3][1], corners[3][2], normal[0], normal[1], normal[2],
		corners[0][0], corners[0][1], corners[0][2], normal[0], normal[1], normal[2],
	)
}

func GenerateCubeVertices(size float32) []float32 {
	half := size / 2
	var vertices []float32

	// Define all 8 corners of the cube
	p := [8][3]float32{
		{-half, -half, -half}, // 0
		{+half, -half, -half}, // 1
		{+half, +half, -half}, // 2
		{-half, +half, -half}, // 3
		{-half, -half, +half}, // 4
		{+half, -half, +half}, // 5
		{+half, +half, +half}, // 6
		{-half, +half, +half}, // 7
	}

	addFace(&vertices, [3]float32{0, 0, 1}, [4][3]float32{p[4], p[5], p[6], p[7]})  // front
	addFace(&vertices, [3]float32{0, 0, -1}, [4][3]float32{p[1], p[0], p[3], p[2]}) // back
	addFace(&vertices, [3]float32{-1, 0, 0}, [4][3]float32{p[0], p[4], p[7], p[3]}) // left
	addFace(&vertices, [3]float32{1, 0, 0}, [4][3]float32{p[5], p[1], p[2], p[6]})  // right
	addFace(&vertices, [3]float32{0, 1, 0}, [4][3]float32{p[3], p[7], p[6], p[2]})  // top
	addFace(&vertices, [3]float32{0, -1, 0}, [4][3]float32{p[0], p[1], p[5], p[4]}) // bottom

	return vertices
}

var wireCubeVertices = []float32{
	// Bottom square
	-1, -1, -1, 1, -1, -1,
	1, -1, -1, 1, -1, 1,
	1, -1, 1, -1, -1, 1,
	-1, -1, 1, -1, -1, -1,

	// Top square
	-1, 1, -1, 1, 1, -1,
	1, 1, -1, 1, 1, 1,
	1, 1, 1, -1, 1, 1,
	-1, 1, 1, -1, 1, -1,

	// Vertical lines
	-1, -1, -1, -1, 1, -1,
	1, -1, -1, 1, 1, -1,
	1, -1, 1, 1, 1, 1,
	-1, -1, 1, -1, 1, 1,
}
