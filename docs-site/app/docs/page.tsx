import { DocPagination } from "@/components/doc-pagination"

export default function DocsPage() {
  return (
    <div className="space-y-6">
      <h1>Introduction</h1>
      <p>
        GoFrame is a Laravel-inspired web framework for Go that provides a clean architecture and essential features for
        building web applications.
      </p>
      <h2>What is GoFrame?</h2>
      <p>
        GoFrame is a web application framework with expressive, elegant syntax. We've already laid the foundation —
        freeing you to create without sweating the small things.
      </p>
      <p>GoFrame takes the pain out of development by easing common tasks used in many web projects, such as:</p>
      <ul>
        <li>Simple, fast routing engine</li>
        <li>Powerful dependency injection container</li>
        <li>Database abstraction with ORM</li>
        <li>Schema migrations</li>
        <li>Robust background job processing</li>
        <li>Real-time event broadcasting</li>
      </ul>
      <h2>Why GoFrame?</h2>
      <p>
        There are a variety of tools and frameworks available to you when building a web application. However, we
        believe GoFrame provides the perfect balance between simplicity and power, giving you a wonderful development
        experience without sacrificing the power that Go provides.
      </p>
      <h3>For Go Developers</h3>
      <p>
        If you're a Go developer looking for a more structured approach to building web applications, GoFrame provides a
        familiar Laravel-like structure while maintaining Go's performance and simplicity.
      </p>
      <h3>For Laravel Developers</h3>
      <p>
        If you're coming from the Laravel ecosystem and want to explore Go, GoFrame provides a familiar structure and
        concepts, making the transition smoother.
      </p>
      <h2>Getting Started</h2>
      <p>
        Ready to start building with GoFrame? Check out our <a href="/docs/installation">installation guide</a> to get
        up and running quickly.
      </p>
      <DocPagination
        next={{
          title: "Installation",
          href: "/docs/installation",
        }}
      />
    </div>
  )
}

