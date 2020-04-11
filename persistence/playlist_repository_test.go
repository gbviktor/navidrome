package persistence

import (
	"github.com/astaxie/beego/orm"
	"github.com/deluan/navidrome/log"
	"github.com/deluan/navidrome/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PlaylistRepository", func() {
	var repo model.PlaylistRepository

	BeforeEach(func() {
		repo = NewPlaylistRepository(log.NewContext(nil), orm.NewOrm())
	})

	Describe("Count", func() {
		It("returns the number of playlists in the DB", func() {
			Expect(repo.CountAll()).To(Equal(int64(2)))
		})
	})

	Describe("Exists", func() {
		It("returns true for an existing playlist", func() {
			Expect(repo.Exists("11")).To(BeTrue())
		})
		It("returns false for a non-existing playlist", func() {
			Expect(repo.Exists("666")).To(BeFalse())
		})
	})

	Describe("Get", func() {
		It("returns an existing playlist", func() {
			p, err := repo.Get("10")
			Expect(err).To(BeNil())
			// Compare all but Tracks and timestamps
			p2 := *p
			p2.Tracks = plsBest.Tracks
			p2.UpdatedAt = plsBest.UpdatedAt
			p2.CreatedAt = plsBest.CreatedAt
			Expect(p2).To(Equal(plsBest))
			// Compare tracks
			for i := range p.Tracks {
				Expect(p.Tracks[i].ID).To(Equal(plsBest.Tracks[i].ID))
			}
		})
		It("returns ErrNotFound for a non-existing playlist", func() {
			_, err := repo.Get("666")
			Expect(err).To(MatchError(model.ErrNotFound))
		})
	})

	Describe("Put/Exists/Delete", func() {
		var newPls model.Playlist
		BeforeEach(func() {
			newPls = model.Playlist{ID: "22", Name: "Great!", Tracks: model.MediaFiles{{ID: "4"}}}
		})
		It("saves the playlist to the DB", func() {
			Expect(repo.Put(&newPls)).To(BeNil())
		})
		It("returns the newly created playlist", func() {
			Expect(repo.Exists("22")).To(BeTrue())
		})
		It("returns deletes the playlist", func() {
			Expect(repo.Delete("22")).To(BeNil())
		})
		It("returns error if tries to retrieve the deleted playlist", func() {
			Expect(repo.Exists("22")).To(BeFalse())
		})
	})

	Describe("GetWithTracks", func() {
		It("returns an existing playlist", func() {
			pls, err := repo.GetWithTracks("10")
			Expect(err).To(BeNil())
			Expect(pls.Name).To(Equal(plsBest.Name))
			Expect(pls.Tracks).To(Equal(model.MediaFiles{
				songDayInALife,
				songRadioactivity,
			}))
		})
	})

	Describe("GetAll", func() {
		It("returns all playlists from DB", func() {
			all, err := repo.GetAll()
			Expect(err).To(BeNil())
			Expect(all[0].ID).To(Equal(plsBest.ID))
			Expect(all[1].ID).To(Equal(plsCool.ID))
		})
	})
})
