// Copyright (C) 2025 T-Force I/O
// This file is part of TFprism
//
// TFprism is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// TFprism is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with TFprism. If not, see <https://www.gnu.org/licenses/>.

package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Link represents a hyperlink to a website.
type Link struct {
	Id uuid.UUID `gorm:"type:uuid;primaryKey;column:id"`

	URL   string `gorm:"column:url;size:512;not null"`
	Title string `gorm:"column:title;size:512;not null"`
}

// Return table name for GORM.
func (Link) TableName() string { return "links" }

// Return a link by ID.
func (c *DbContext) GetLink(id uuid.UUID) (*Link, error) {
	return c.findLinkByID(c.db, id)
}

// Return a link by URL.
func (c *DbContext) GetLinkByURL(url string) (*Link, error) {
	return c.findLinkByURL(c.db, url)
}

// Return links by ID. Returns all links when ids is empty.
func (c *DbContext) GetLinks(ids []uuid.UUID) ([]*Link, error) {
	if len(ids) == 0 {
		return c.findLinks(c.db)
	}
	return c.findLinksByID(c.db, ids)
}

// Return links by URL.
func (c *DbContext) GetLinksByURL(urls []string) ([]*Link, error) {
	return c.findLinksByURL(c.db, urls)
}

// Insert a link to database.
func (c *DbContext) InsertLink(link *Link) error {
	tx := c.db.Begin()
	err := c.insertLink(tx, link)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// Insert links to database.
func (c *DbContext) InsertLinks(links []*Link) error {
	tx := c.db.Begin()
	err := c.insertLinks(tx, links)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// Update a link in database.
func (c *DbContext) UpdateLink(id uuid.UUID, link *Link) error {
	tx := c.db.Begin()
	err := c.updateLinkById(tx, id, link)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// Delete a link by ID.
func (c *DbContext) DeleteLink(id uuid.UUID) error {
	tx := c.db.Begin()
	err := c.deleteLinkById(tx, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// Return a link with given ID from transaction.
func (c *DbContext) findLinkByID(tx *gorm.DB, id uuid.UUID) (*Link, error) {
	var doc *Link
	result := tx.Model(&Link{}).
		Where("id = ?", id).
		First(&doc)
	if c.isEmptyResultError(result.Error) {
		return nil, nil
	}
	return doc, result.Error
}

// Return a link with given URL from transaction.
func (c *DbContext) findLinkByURL(tx *gorm.DB, url string) (*Link, error) {
	var doc *Link
	result := tx.Model(&Link{}).
		Where("url = ?", url).
		First(&doc)
	if c.isEmptyResultError(result.Error) {
		return nil, nil
	}
	return doc, result.Error
}

// Return all links.
func (c *DbContext) findLinks(tx *gorm.DB) ([]*Link, error) {
	var docs []*Link
	result := tx.Model(&Link{}).
		Find(&docs)
	return docs, result.Error
}

// Return links with given IDs from transaction.
func (c *DbContext) findLinksByID(tx *gorm.DB, ids []uuid.UUID) ([]*Link, error) {
	var docs []*Link
	result := tx.Model(&Link{}).
		Where("id IN ?", ids).
		Find(&docs)
	return docs, result.Error
}

// Return links with given URLs from transaction.
func (c *DbContext) findLinksByURL(tx *gorm.DB, urls []string) ([]*Link, error) {
	var docs []*Link
	result := tx.Model(&Link{}).
		Where("url IN ?", urls).
		Find(&docs)
	return docs, result.Error
}

// Insert a link into database using given transaction.
func (c *DbContext) insertLink(tx *gorm.DB, link *Link) error {
	link.Id = c.newRandomID(tx, RefTypeLink)
	result := tx.Create(link)
	return result.Error
}

// Insert links into database using given transaction.
func (c *DbContext) insertLinks(tx *gorm.DB, links []*Link) error {
	for _, link := range links {
		link.Id = c.newRandomID(tx, RefTypeLink)
		result := tx.Create(link)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

// Update a link with given ID in database using given transaction.
func (c *DbContext) updateLinkById(tx *gorm.DB, id uuid.UUID, link *Link) error {
	result := tx.Model(&Link{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"url":   link.URL,
			"title": link.Title,
		})
	return result.Error
}

// Delete a link with given ID from database using given transaction.
func (c *DbContext) deleteLinkById(tx *gorm.DB, id uuid.UUID) error {
	result := tx.Where("id = ?", id).
		Delete(&Link{})
	return result.Error
}
