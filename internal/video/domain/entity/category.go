package entity

type Category struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    ParentID  int    `json:"parent_id"`
    SortOrder int    `json:"sort_order"`
}
