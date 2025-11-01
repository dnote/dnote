/* Copyright 2025 Dnote Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package database

import (
	"time"
)

// Model is the base model definition
type Model struct {
	ID        int       `gorm:"primaryKey" json:"-"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Book is a model for a book
type Book struct {
	Model
	UUID      string `json:"uuid" gorm:"uniqueIndex;type:text"`
	UserID    int    `json:"user_id" gorm:"index"`
	Label     string `json:"label" gorm:"index"`
	Notes     []Note `json:"notes" gorm:"foreignKey:BookUUID;references:UUID"`
	AddedOn   int64  `json:"added_on"`
	EditedOn  int64  `json:"edited_on"`
	USN       int    `json:"-" gorm:"index"`
	Deleted   bool   `json:"-" gorm:"default:false"`
}

// Note is a model for a note
type Note struct {
	Model
	UUID      string `json:"uuid" gorm:"index;type:text"`
	Book      Book   `json:"book" gorm:"foreignKey:BookUUID;references:UUID"`
	User      User   `json:"user"`
	UserID    int    `json:"user_id" gorm:"index"`
	BookUUID  string `json:"book_uuid" gorm:"index;type:text"`
	Body      string `json:"content"`
	AddedOn   int64  `json:"added_on"`
	EditedOn  int64  `json:"edited_on"`
	USN       int    `json:"-" gorm:"index"`
	Deleted   bool   `json:"-" gorm:"default:false"`
	Client    string `gorm:"index"`
}

// User is a model for a user
type User struct {
	Model
	UUID           string     `json:"uuid" gorm:"type:text;index"`
	Email          NullString `gorm:"index"`
	Password       NullString `json:"-"`
	LastLoginAt    *time.Time `json:"-"`
	MaxUSN         int        `json:"-" gorm:"default:0"`
	FullSyncBefore int64      `json:"-" gorm:"default:0"`
}

// Token is a model for a token
type Token struct {
	Model
	UserID int    `gorm:"index"`
	Value  string `gorm:"index"`
	Type   string
	UsedAt *time.Time
}

// Session represents a user session
type Session struct {
	Model
	UserID     int    `gorm:"index"`
	Key        string `gorm:"index"`
	LastUsedAt time.Time
	ExpiresAt  time.Time
}
