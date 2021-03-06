package tablecloth

import (
	"net/http"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestManager(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Manager")
}

var _ = Describe("Adding servers", func() {
	var (
		setupCount int
	)

	BeforeEach(func() {
		theManager = &manager{}
		setupCount = 0
		setupFunc = func() {
			theManager.servers = make(map[string]*serverInfo)
			setupCount++
		}
	})

	AfterEach(func() {
		for _, si := range theManager.servers {
			si.server.Close()
		}
	})

	It("Should add the listener using the given ident", func() {
		go ListenAndServe("127.0.0.1:8081", http.NotFoundHandler(), "one")
		time.Sleep(10 * time.Millisecond)

		listener := theManager.servers["one"].listener
		Expect(listener.Addr().String()).To(Equal("127.0.0.1:8081"))
	})

	It("Should use an ident of 'default' if none given", func() {
		go ListenAndServe("127.0.0.1:8081", http.NotFoundHandler())
		time.Sleep(10 * time.Millisecond)

		listener := theManager.servers["default"].listener
		Expect(listener.Addr().String()).To(Equal("127.0.0.1:8081"))
	})

	Context("listening on multiple addresses", func() {
		It("Should allow listening on multiple addresses", func() {
			go ListenAndServe("127.0.0.1:8081", http.NotFoundHandler(), "one")
			go ListenAndServe("127.0.0.1:8082", http.NotFoundHandler(), "two")
			time.Sleep(10 * time.Millisecond)

			listener := theManager.servers["one"].listener
			Expect(listener.Addr().String()).To(Equal("127.0.0.1:8081"))

			listener = theManager.servers["two"].listener
			Expect(listener.Addr().String()).To(Equal("127.0.0.1:8082"))
		})

		It("Should only run the setup function once", func() {
			go ListenAndServe("127.0.0.1:8081", http.NotFoundHandler(), "one")
			go ListenAndServe("127.0.0.1:8081", http.NotFoundHandler(), "two")
			time.Sleep(10 * time.Millisecond)

			Expect(setupCount).To(Equal(1))
		})

		It("Should return an error if given duplicate idents", func() {
			go ListenAndServe("127.0.0.1:8081", http.NotFoundHandler(), "foo")
			time.Sleep(10 * time.Millisecond)
			err := ListenAndServe("127.0.0.1:8082", http.NotFoundHandler(), "foo")

			Expect(err).To(MatchError("duplicate ident"))

			listener := theManager.servers["foo"].listener
			Expect(listener.Addr().String()).To(Equal("127.0.0.1:8081"))
		})
	})
})
