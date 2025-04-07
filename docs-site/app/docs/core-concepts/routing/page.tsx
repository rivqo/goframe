import { DocPagination } from "@/components/doc-pagination"

export default function RoutingPage() {
  return (
    <div className="space-y-6">
      <h1>Routing</h1>
      <p>
        GoFrame provides a simple and expressive routing system that makes it easy to define routes for your web
        application. The router supports various HTTP methods, route groups, middleware, and more.
      </p>

      <h2>Basic Routing</h2>
      <p>
        The most basic routes accept a URI and a closure, providing a simple and expressive method of defining routes:
      </p>
      <pre>
        <code>{`// routes/web.go
r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, World!"))
})

r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
    // Create a new user
})

r.Put("/users/:id", func(w http.ResponseWriter, r *http.Request) {
    // Update a user
})

r.Delete("/users/:id", func(w http.ResponseWriter, r *http.Request) {
    // Delete a user
})`}</code>
      </pre>

      <h2>Route Parameters</h2>
      <p>
        Sometimes you need to capture segments of the URI within your route. For example, you may need to capture a
        user's ID from the URL:
      </p>
      <pre>
        <code>{`r.Get("/users/:id", func(w http.ResponseWriter, r *http.Request) {
    // Get the ID from the URL
    id := router.GetParam(r, "id")
    
    // Use the ID to fetch the user
    user, err := userRepo.FindByID(id)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    
    // Return the user
    json.NewEncoder(w).Encode(user)
})`}</code>
      </pre>

      <h2>Route Groups</h2>
      <p>
        Route groups allow you to share route attributes, such as middleware or prefixes, across a large number of
        routes without needing to define those attributes on each individual route:
      </p>
      <pre>
        <code>{`// Create a group with a prefix
api := r.Group("/api")

// Add middleware to the group
api.Use(middleware.JSON())
api.Use(middleware.CORS())

// Define routes for the group
api.Get("/users", userController.Index)
api.Get("/users/:id", userController.Show)
api.Post("/users", userController.Store)
api.Put("/users/:id", userController.Update)
api.Delete("/users/:id", userController.Destroy)`}</code>
      </pre>

      <h2>Named Routes</h2>
      <p>Named routes allow you to generate URLs to specific routes using a name instead of the full URL:</p>
      <pre>
        <code>{`// Define a named route
r.Get("/users/:id", userController.Show).Name("users.show")

// Generate a URL to the named route
url := router.URL("users.show", map[string]string{"id": "1"})
// url = "/users/1"`}</code>
      </pre>

      <h2>Route Middleware</h2>
      <p>Middleware provide a convenient mechanism for filtering HTTP requests entering your application:</p>
      <pre>
        <code>{`// Apply middleware to a single route
r.Get("/admin", adminController.Index).Middleware(middleware.Auth())

// Apply middleware to a group of routes
admin := r.Group("/admin")
admin.Use(middleware.Auth())
admin.Use(middleware.Role("admin"))

admin.Get("/dashboard", adminController.Dashboard)
admin.Get("/users", adminController.Users)`}</code>
      </pre>

      <h2>Static Files</h2>
      <p>GoFrame makes it easy to serve static files like images, CSS, and JavaScript:</p>
      <pre>
        <code>{`// Serve a directory of static files
r.Static("/assets", "./public/assets")

// Serve a single static file
r.StaticFile("/favicon.ico", "./public/favicon.ico")`}</code>
      </pre>

      <h2>404 Handler</h2>
      <p>You can define a custom handler for 404 errors:</p>
      <pre>
        <code>{`// Set a custom 404 handler
r.NotFound(func(w http.ResponseWriter, r *http.Request) {
    view.Render(w, "errors/404", map[string]interface{}{
        "title": "Page Not Found",
    })
})`}</code>
      </pre>

      <DocPagination
        prev={{
          title: "Core Concepts",
          href: "/docs/core-concepts",
        }}
        next={{
          title: "Controllers",
          href: "/docs/core-concepts/controllers",
        }}
      />
    </div>
  )
}

