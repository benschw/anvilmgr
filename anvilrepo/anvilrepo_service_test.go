package anvilrepo

import (
	"fmt"
	"git.bvops.net/scm/auto/anvilmgr.git/anvilrepo/api"
	"git.bvops.net/scm/auto/anvilmgr.git/anvilrepo/client"
	"github.com/benschw/opin-go/rando"
	"github.com/elazarl/go-bindata-assetfs"
	. "gopkg.in/check.v1"
	"io/ioutil"
	"log"
	"syscall"
)

var _ = fmt.Print
var _ = log.Print
var _ = syscall.Unlink

type TestSuite struct {
	s    *AnvilRepoService
	path string
	host string
}

var _ = Suite(&TestSuite{})

func (s *TestSuite) SetUpSuite(c *C) {

	host := fmt.Sprintf("localhost:%d", rando.Port())

	f, err := ioutil.TempDir("", "anviltest")
	if err != nil {
		panic(err)
	}
	s.path = f

	server, err := NewAnvilRepoService(host, s.path, &assetfs.AssetFS{})
	if err != nil {
		panic(err)
	}
	s.s = server
	go s.s.Run()

	s.host = "http://" + host
}

func (s *TestSuite) SetUpTest(c *C) {
}

func (s *TestSuite) TearDownTest(c *C) {
	paths, err := ioutil.ReadDir(s.path)
	if err != nil {
		panic(err)
	}

	for _, path := range paths {
		syscall.Unlink(path.Name())
	}
}
func (s *TestSuite) TearDownSuite(c *C) {
	syscall.Unlink(s.path)
}

// Location should be added
func (s *TestSuite) TestAdd(c *C) {
	// given
	expected, _ := api.NewArtifact("puppetlabs", "stdlib", "4.5.1", "puppetlabs-stdlib-4.5.1.tar.gz")

	cl := &client.AnvilRepoClient{Host: s.host}

	// when
	created, err := cl.AddArtifact("puppetlabs", "stdlib", "../puppetlabs-stdlib-4.5.1.tar.gz")

	// then
	c.Assert(err, Equals, nil)

	c.Assert(expected, DeepEquals, &created)
}

func (s *TestSuite) TestFindAll(c *C) {
	// given
	expected, _ := api.NewArtifact("puppetlabs", "stdlib", "4.5.1", "puppetlabs-stdlib-4.5.1.tar.gz")

	cl := &client.AnvilRepoClient{Host: s.host}
	cl.AddArtifact("puppetlabs", "stdlib", "../puppetlabs-stdlib-4.5.1.tar.gz")
	cl.AddArtifact("foo", "bar", "../puppetlabs-stdlib-4.5.1.tar.gz")

	// when
	found, err := cl.FindAllArtifacts()
	// then
	c.Assert(err, Equals, nil)

	c.Assert(expected, DeepEquals, &found[1])
	c.Assert(len(found), Equals, 2)
}

func (s *TestSuite) TestFindAllVersions(c *C) {
	// given
	expected, _ := api.NewArtifact("puppetlabs", "stdlib", "4.5.1", "puppetlabs-stdlib-4.5.1.tar.gz")

	cl := &client.AnvilRepoClient{Host: s.host}
	cl.AddArtifact("puppetlabs", "stdlib", "../puppetlabs-stdlib-4.5.1.tar.gz")
	cl.AddArtifact("foo", "bar", "../puppetlabs-stdlib-4.5.1.tar.gz")

	// when
	found, err := cl.FindAllArtifactVersions("puppetlabs", "stdlib")
	// then
	c.Assert(err, Equals, nil)

	c.Assert(expected, DeepEquals, &found[0])
	c.Assert(len(found), Equals, 1)
}

// // Client should return ErrStatusBadRequest when entity doesn't validate
// func (s *TestSuite) TestAddBadRequest(c *C) {
// 	// given
// 	locClient := client.LocationClient{Host: s.host}

// 	// when
// 	_, err := locClient.AddLocation("", "Texas")

// 	// then
// 	c.Assert(err, Equals, rest.ErrStatusBadRequest)
// }

// // Client should return ErrStatusConflict when id exists
// // (not supported by client so pulled impl into test)
// func (s *TestSuite) TestAddConflict(c *C) {
// 	// given
// 	locClient := client.LocationClient{Host: s.host}
// 	created, _ := locClient.AddLocation("Austin", "Texas")

// 	// when
// 	url := fmt.Sprintf("%s/location", s.host)
// 	r, _ := rest.MakeRequest("POST", url, created)
// 	err := rest.ProcessResponseEntity(r, nil, http.StatusCreated)

// 	// then
// 	c.Assert(err, Equals, rest.ErrStatusConflict)
// }

// // Location should be findable
// func (s *TestSuite) TestFind(c *C) {
// 	// given
// 	locClient := client.LocationClient{Host: s.host}
// 	created, _ := locClient.AddLocation("Austin", "Texas")

// 	// when
// 	found, err := locClient.FindLocation(created.Id)

// 	// then
// 	c.Assert(err, Equals, nil)

// 	c.Assert(created, DeepEquals, found)
// }

// // Client should return ErrStatusNotFound when not found
// func (s *TestSuite) TestFindNotFound(c *C) {
// 	// given
// 	locClient := client.LocationClient{Host: s.host}

// 	// when
// 	_, err := locClient.FindLocation(1)

// 	// then
// 	c.Assert(err, Equals, rest.ErrStatusNotFound)
// }

// // Client should return ErrStatusBadRequest when id doesn't validate
// // (not supported by client so pulled impl into test)
// func (s *TestSuite) TestFindBadRequest(c *C) {

// 	// when
// 	url := fmt.Sprintf("%s/location/%s", s.host, "asd")
// 	r, err := rest.MakeRequest("GET", url, nil)
// 	err = rest.ProcessResponseEntity(r, nil, http.StatusOK)

// 	// then
// 	c.Assert(err, Equals, rest.ErrStatusBadRequest)
// }

// // Find all should return all locations
// func (s *TestSuite) TestFindAll(c *C) {
// 	// given
// 	locClient := client.LocationClient{Host: s.host}

// 	loc1, err := locClient.AddLocation("Austin", "Texas")
// 	loc2, err := locClient.AddLocation("Williamsburg", "Virginia")
// 	// when

// 	foundLocations, err := locClient.FindAllLocations()

// 	// then
// 	c.Assert(err, Equals, nil)

// 	c.Assert(foundLocations, DeepEquals, []api.Location{loc1, loc2})
// }

// // Find all should return empty list when no results are found
// func (s *TestSuite) TestFindAllEmpty(c *C) {
// 	// given
// 	locClient := client.LocationClient{Host: s.host}

// 	// when
// 	foundLocations, err := locClient.FindAllLocations()

// 	// then
// 	c.Assert(err, Equals, nil)

// 	c.Assert(len(foundLocations), Equals, 0)
// }

// // Save should update a location
// func (s *TestSuite) TestSave(c *C) {
// 	// given
// 	locClient := client.LocationClient{Host: s.host}

// 	location, _ := locClient.AddLocation("Austin", "Texas")

// 	// when
// 	saved, err := locClient.SaveLocation(location)

// 	// then
// 	c.Assert(err, Equals, nil)

// 	c.Assert(location.State, DeepEquals, saved.State)
// }

// // Client should return ErrStatusNotFound if trying to save to an id that doesn't exist
// func (s *TestSuite) TestSaveNotFound(c *C) {
// 	// given
// 	locClient := client.LocationClient{Host: s.host}

// 	location, _ := locClient.AddLocation("Austin", "Texas")

// 	// when
// 	location.Id = location.Id + 1
// 	location.State = "foo"
// 	_, err := locClient.SaveLocation(location)

// 	// then
// 	c.Assert(err, Equals, rest.ErrStatusNotFound)
// }

// // Client should return ErrStatusBadRequest if entity doesn't validate
// func (s *TestSuite) TestSaveBadRequestFromEntity(c *C) {
// 	// given
// 	locClient := client.LocationClient{Host: s.host}

// 	location, _ := locClient.AddLocation("Austin", "Texas")

// 	// when
// 	location.State = ""
// 	_, err := locClient.SaveLocation(location)

// 	// then
// 	c.Assert(err, Equals, rest.ErrStatusBadRequest)
// }

// // Client should return ErrStatusBadRequest if Id doesn't validate
// // (not supported by client so pulled impl into test)
// func (s *TestSuite) TestSaveBadRequestFromId(c *C) {
// 	// given
// 	locClient := client.LocationClient{Host: s.host}

// 	location, _ := locClient.AddLocation("Austin", "Texas")

// 	// when
// 	url := fmt.Sprintf("%s/location/%s", s.host, "asd")
// 	r, err := rest.MakeRequest("GET", url, location)
// 	err = rest.ProcessResponseEntity(r, nil, http.StatusOK)

// 	// then
// 	c.Assert(err, Equals, rest.ErrStatusBadRequest)
// }

// // Delete should Delete a location
// func (s *TestSuite) TestDelete(c *C) {
// 	// given
// 	locClient := client.LocationClient{Host: s.host}

// 	location, _ := locClient.AddLocation("Austin", "Texas")

// 	// when
// 	err := locClient.DeleteLocation(location.Id)

// 	// then
// 	c.Assert(err, Equals, nil)

// 	foundLocations, _ := locClient.FindAllLocations()

// 	c.Assert(len(foundLocations), Equals, 0)
// }

// // Client should return ErrStatusNotFound if trying to delete an Id that doesn't exist
// func (s *TestSuite) TestDeleteNotFound(c *C) {
// 	// given
// 	locClient := client.LocationClient{Host: s.host}

// 	// when
// 	err := locClient.DeleteLocation(1)

// 	// then
// 	c.Assert(err, Equals, rest.ErrStatusNotFound)
// }

// // Client should return ErrStatusBadRequest if Id doesn't validate
// // (not supported by client so pulled impl into test)
// func (s *TestSuite) TestDeleteBdRequesta(c *C) {
// 	// when
// 	url := fmt.Sprintf("%s/location/%s", s.host, "asd")
// 	r, _ := rest.MakeRequest("DELETE", url, nil)
// 	err := rest.ProcessResponseEntity(r, nil, http.StatusNoContent)

// 	// then
// 	c.Assert(err, Equals, rest.ErrStatusBadRequest)
// }
