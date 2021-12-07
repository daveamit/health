package health

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEnsureItem(t *testing.T) {

	hStruct := defaultHealth.(*healthImpl)
	EnsureService("redis", "namespace_0")

	if len(hStruct.items) != 1 {
		t.Errorf("items must have just one item")
		return
	}

	i := hStruct.items[0]
	if i.Name != "redis" || i.Namespace != "namespace_0" {
		t.Errorf("name should be redis and namespace should be namespace_0")
		return
	}

	EnsureService("redis", "namespace_0")

	if len(hStruct.items) != 1 {
		t.Errorf("items must have just one item (ensure should not add new item)")
		return
	}

	i = hStruct.items[0]
	if i.Name != "redis" || i.Namespace != "namespace_0" {
		t.Errorf("name should be redis and namespace should be namespace_0")
		return
	}

	EnsureService("redis", "namespace_1")

	if len(hStruct.items) != 2 {
		t.Errorf("items must have just two item")
		return
	}

	i = hStruct.items[0]
	if i.Name != "redis" || i.Namespace != "namespace_0" {
		t.Errorf("first element should be: redis, namespace_0")
		return
	}
	i = hStruct.items[1]
	if i.Name != "redis" || i.Namespace != "namespace_1" {
		t.Errorf("second element should be: redis, namespace_1")
		return
	}

}

func TestHealthCheckHandler(t *testing.T) {
	EnsureService("redis", "namespace_0")

	// MUST return 500 if ServiceUp or ServiceDown was never called
	{ // Create a request to pass to our handler. We don't have any query parameters for now, so we'll
		// pass 'nil' as the third parameter.
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := HealthCheckHandler()

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check the response body is what we expect.
		expected := `[{"name":"redis","namespace":"namespace_0","state":"undefined"},{"name":"redis","namespace":"namespace_1","state":"undefined"}]`
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	}
	// MUST return 200 if ServiceUp was called
	{ // Create a request to pass to our handler. We don't have any query parameters for now, so we'll
		ServiceUp("redis", "namespace_0")
		ServiceUp("redis", "namespace_1")

		// pass 'nil' as the third parameter.
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := HealthCheckHandler()

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check the response body is what we expect.
		expected := `[{"name":"redis","namespace":"namespace_0","state":"Up"},{"name":"redis","namespace":"namespace_1","state":"Up"}]`
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	}

	// MUST return 400 if ServiceDown was called
	{ // Create a request to pass to our handler. We don't have any query parameters for now, so we'll
		ServiceDown("redis", "namespace_0")

		// pass 'nil' as the third parameter.
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := HealthCheckHandler()

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check the response body is what we expect.
		expected := `[{"name":"redis","namespace":"namespace_0","state":"Down"},{"name":"redis","namespace":"namespace_1","state":"Up"}]`
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	}

	EnsureService("redis", "namespace_1")
	// MUST return 500 if ServiceUp or ServiceDown was never called
	{ // Create a request to pass to our handler. We don't have any query parameters for now, so we'll
		// pass 'nil' as the third parameter.
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := HealthCheckHandler()

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check the response body is what we expect.
		expected := `[{"name":"redis","namespace":"namespace_0","state":"Down"},{"name":"redis","namespace":"namespace_1","state":"Up"}]`
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	}

	// MUST return 200 if ServiceUp was called (upping downed service)
	{ // Create a request to pass to our handler. We don't have any query parameters for now, so we'll
		ServiceUp("redis", "namespace_0")
		ServiceUp("redis", "namespace_1")

		// pass 'nil' as the third parameter.
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := HealthCheckHandler()

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check the response body is what we expect.
		expected := `[{"name":"redis","namespace":"namespace_0","state":"Up"},{"name":"redis","namespace":"namespace_1","state":"Up"}]`
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	}

	// MUST return 400 if ServiceUp was called on any of the services (downing uped service)
	{ // Create a request to pass to our handler. We don't have any query parameters for now, so we'll
		ServiceDown("redis", "namespace_1")

		// pass 'nil' as the third parameter.
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := HealthCheckHandler()

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check the response body is what we expect.
		expected := `[{"name":"redis","namespace":"namespace_0","state":"Up"},{"name":"redis","namespace":"namespace_1","state":"Down"}]`
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	}
}
