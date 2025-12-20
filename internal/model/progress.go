package model

type ReadingProgress struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	ItemType  string  `json:"item_type"` // "writer" или "article"
	ItemID    int     `json:"item_id"`
	Progress  float64 `json:"progress"` // 0.0 - 1.0
	UpdatedAt string  `json:"updated_at"`
}

type BulkProgressRequest struct {
	Writers  map[string]float64 `json:"writers"`  // "101": 0.8
	Articles map[string]float64 `json:"articles"` // "201": 0.5
}

type ProgressUpdateRequest struct {
	Type     string  `json:"type" binding:"required,oneof=writer article"`
	ID       int     `json:"id" binding:"required,min=1"`
	Progress float64 `json:"progress" binding:"required,min=0,max=1"`
}

type ProgressResponse struct {
	Writers  map[string]float64 `json:"writers"`  // "101": 0.8
	Articles map[string]float64 `json:"articles"` // "201": 0.5
}
