/*
* victi.ms API microservice
* Copyright (C) 2017 The victi.ms team
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU Affero General Public License as published
* by the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU Affero General Public License for more details.
*
* You should have received a copy of the GNU Affero General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package api

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/victims/victims-common/log"
	"github.com/victims/victims-common/types"
)

// Upload checks a provided package against the victims database
func Upload(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("package")
	if err != nil {
		log.Logger.Infof("Error getting the file in upload: %s", err)
		// Bad Request
		c.AbortWithStatus(400)
	}
	defer file.Close()

	// Write the file out to the file system
	// TODO: Don't use the original name or /tmp
	out, err := os.Create("/tmp/" + fileHeader.Filename + ".test")
	defer out.Close()
	if err != nil {
		log.Logger.Fatal(err)
	}
	// Copy the content from the upload file to the file on disk
	fileSize, err := io.Copy(out, file)
	if err != nil {
		log.Logger.Fatal(err)
	}
	log.Logger.Infof("%s: %vk", fileHeader.Filename, int64(fileSize/1024))
	// TODO: Submit to hashing service
	// TODO: Retrieve response and store in requestedHash
	requestedHash := types.MultipleHashRequest{}
	requestedHash.Hashes = append(requestedHash.Hashes, types.SingleHashRequest{
		Hash: "a0a86214ea153fb07ff35ceec0848dd1703eae22de036a825efc8394e50f65e3044832f3b49cf7e45a39edc470bdf738abc36a3a78ca7df3a6e73c14eaef94a8",
	})

	// Fall through to a 404
	c.AbortWithStatus(404)
}

// HashMounts mounts all hash related routes to a router
func HashMounts(router *gin.Engine) {
	router.POST("/upload", Upload)
}
