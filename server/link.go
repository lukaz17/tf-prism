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

package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tforceaio/tf-prism/db"
)

// linkRequest represents the request body for creating and updating links.
type linkRequest struct {
	URL   string `json:"url" binding:"required"`
	Title string `json:"title" binding:"required"`
}

// Register link management routes on the router.
func registerLinkRoutes(r *gin.Engine, ctx *db.DbContext) {
	r.GET("/links", func(c *gin.Context) {
		links, err := ctx.GetLinks(nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, links)
	})

	r.GET("/links/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid link id"})
			return
		}
		link, err := ctx.GetLink(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if link == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "link not found"})
			return
		}
		c.JSON(http.StatusOK, link)
	})

	r.POST("/links", func(c *gin.Context) {
		var req linkRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		link := &db.Link{URL: req.URL, Title: req.Title}
		if err := ctx.InsertLink(link); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, link)
	})

	r.PUT("/links/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid link id"})
			return
		}
		var req linkRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := ctx.UpdateLink(id, &db.Link{URL: req.URL, Title: req.Title}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	r.DELETE("/links/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid link id"})
			return
		}
		if err := ctx.DeleteLink(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	})
}
