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
	"github.com/victims/victims-common/db"
	"github.com/victims/victims-upload/upload"
	"net/http"
	"time"
	"io/ioutil"
	"encoding/json"
)

// Upload looks up hash from services and persists hash to victims database
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
	out, err := os.Create("/tmp/" + fileHeader.Filename)
	defer out.Close()
	if err != nil {
		log.Logger.Fatal(err)
	}
	// Copy the content from the upload file to the file on disk
	len, err := io.Copy(out, file)
	if err != nil {
		log.Logger.Fatal(err)
	}

	log.Logger.Infof("length: %v", len)

	request, err := upload.UploadRequest("http://localhost:8081/hash", "library2", "/tmp/" + fileHeader.Filename)
	if err != nil {
		log.Logger.Error(err)
		c.AbortWithStatus(500)
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Logger.Error(err)
		c.AbortWithStatus(500)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Logger.Error(err)
		c.AbortWithStatus(500)
	}

	multiHash := []types.Hash{}
	if err := json.Unmarshal(body, &multiHash); err != nil {
		log.Logger.Error(err)
		c.AbortWithStatus(500)
	}

	cves := types.CVEs{}
	cve := c.Request.FormValue("cve")
	cves.AppendSingle(cve)

	col, _ := db.GetCollection("hashes")
	for _ , hash := range multiHash {
		hash.Cves = cves
		col.Insert(hash)
	}


	//log.Logger.Infof("length: %v, %v", len, string(responseJson))

	//TODO persist singleHash to DB with CVE and Submitter

	// Fall through to a 404

	c.AbortWithStatus(404)
}


// HashMounts mounts all hash related routes to a router
func HashMounts(router *gin.Engine) {
	router.POST("/upload", Upload)
}
