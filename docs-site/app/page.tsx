import Link from "next/link"
import { ArrowRight } from "lucide-react"

import { siteConfig } from "@/config/site"
import { cn } from "@/lib/utils"
import { buttonVariants } from "@/components/ui/button"
import { SiteHeader } from "@/components/site-header"

export default function Home() {
  return (
    <div className="flex min-h-screen flex-col">
      <SiteHeader />
      <main className="flex-1 mx-auto">
        <section className="space-y-6 pb-8 pt-6 md:pb-12 md:pt-10 lg:py-32">
          <div className="container flex max-w-[64rem] flex-col items-center gap-4 text-center">
            <Link
              href="https://github.com/goframe/goframe"
              className="rounded-2xl bg-muted px-4 py-1.5 text-sm font-medium"
              target="_blank"
            >
              Follow along on GitHub
            </Link>
            <h1 className="text-4xl font-bold tracking-tight sm:text-5xl md:text-6xl lg:text-7xl">GoFrame</h1>
            <p className="max-w-[42rem] leading-normal text-muted-foreground sm:text-xl sm:leading-8">
              A Laravel-inspired web framework for Go that provides a clean architecture and essential features for
              building web applications.
            </p>
            <div className="flex flex-wrap justify-center gap-4">
              <Link href="/docs" className={cn(buttonVariants({ size: "lg" }))}>
                Documentation
              </Link>
              <Link
                href={siteConfig.links.github}
                target="_blank"
                rel="noreferrer"
                className={cn(buttonVariants({ variant: "outline", size: "lg" }))}
              >
                GitHub
              </Link>
            </div>
          </div>
        </section>
        <section className="container space-y-6 py-8 md:py-12 lg:py-24">
          <div className="mx-auto flex max-w-[58rem] flex-col items-center space-y-4 text-center">
            <h2 className="text-3xl font-bold leading-[1.1] sm:text-3xl md:text-5xl">Features</h2>
            <p className="max-w-[85%] leading-normal text-muted-foreground sm:text-lg sm:leading-7">
              GoFrame provides a comprehensive set of features to help you build robust web applications.
            </p>
          </div>
          <div className="mx-auto grid justify-center gap-4 sm:grid-cols-2 md:max-w-[64rem] md:grid-cols-3">
            <div className="relative overflow-hidden rounded-lg border bg-card p-6">
              <div className="flex h-full flex-col justify-between">
                <div className="space-y-2">
                  <h3 className="font-bold">ORM System</h3>
                  <p className="text-sm text-muted-foreground">
                    Simple database abstraction with repositories for clean data access.
                  </p>
                </div>
                <Link href="/docs/database" className={cn(buttonVariants({ variant: "link", size: "sm" }), "px-0")}>
                  Learn more <ArrowRight className="ml-1 h-4 w-4" />
                </Link>
              </div>
            </div>
            <div className="relative overflow-hidden rounded-lg border bg-card p-6">
              <div className="flex h-full flex-col justify-between">
                <div className="space-y-2">
                  <h3 className="font-bold">Authentication</h3>
                  <p className="text-sm text-muted-foreground">
                    JWT-based authentication with middleware for protected routes.
                  </p>
                </div>
                <Link
                  href="/docs/features/authentication"
                  className={cn(buttonVariants({ variant: "link", size: "sm" }), "px-0")}
                >
                  Learn more <ArrowRight className="ml-1 h-4 w-4" />
                </Link>
              </div>
            </div>
            <div className="relative overflow-hidden rounded-lg border bg-card p-6">
              <div className="flex h-full flex-col justify-between">
                <div className="space-y-2">
                  <h3 className="font-bold">Migrations</h3>
                  <p className="text-sm text-muted-foreground">
                    Database schema versioning for easy database management.
                  </p>
                </div>
                <Link
                  href="/docs/database/migrations"
                  className={cn(buttonVariants({ variant: "link", size: "sm" }), "px-0")}
                >
                  Learn more <ArrowRight className="ml-1 h-4 w-4" />
                </Link>
              </div>
            </div>
            <div className="relative overflow-hidden rounded-lg border bg-card p-6">
              <div className="flex h-full flex-col justify-between">
                <div className="space-y-2">
                  <h3 className="font-bold">Resources</h3>
                  <p className="text-sm text-muted-foreground">
                    API resource transformers for consistent API responses.
                  </p>
                </div>
                <Link
                  href="/docs/features/resources"
                  className={cn(buttonVariants({ variant: "link", size: "sm" }), "px-0")}
                >
                  Learn more <ArrowRight className="ml-1 h-4 w-4" />
                </Link>
              </div>
            </div>
            <div className="relative overflow-hidden rounded-lg border bg-card p-6">
              <div className="flex h-full flex-col justify-between">
                <div className="space-y-2">
                  <h3 className="font-bold">Rate Limiting</h3>
                  <p className="text-sm text-muted-foreground">
                    Configurable request rate limiting to protect your application.
                  </p>
                </div>
                <Link
                  href="/docs/features/rate-limiting"
                  className={cn(buttonVariants({ variant: "link", size: "sm" }), "px-0")}
                >
                  Learn more <ArrowRight className="ml-1 h-4 w-4" />
                </Link>
              </div>
            </div>
            <div className="relative overflow-hidden rounded-lg border bg-card p-6">
              <div className="flex h-full flex-col justify-between">
                <div className="space-y-2">
                  <h3 className="font-bold">CLI Commands</h3>
                  <p className="text-sm text-muted-foreground">Code generation and project management tools.</p>
                </div>
                <Link
                  href="/docs/features/cli-commands"
                  className={cn(buttonVariants({ variant: "link", size: "sm" }), "px-0")}
                >
                  Learn more <ArrowRight className="ml-1 h-4 w-4" />
                </Link>
              </div>
            </div>
          </div>
        </section>
        <section className="container py-8 md:py-12 lg:py-24">
          <div className="mx-auto flex max-w-[58rem] flex-col items-center justify-center gap-4 text-center">
            <h2 className="text-3xl font-bold leading-[1.1] sm:text-3xl md:text-5xl">Get Started</h2>
            <p className="max-w-[85%] leading-normal text-muted-foreground sm:text-lg sm:leading-7">
              Start building your Go application with GoFrame today.
            </p>
            <Link href="/docs/installation" className={cn(buttonVariants({ size: "lg" }))}>
              Installation Guide
            </Link>
          </div>
        </section>
      </main>
      <footer className="border-t py-6 md:py-0">
        <div className="container flex flex-col items-center justify-between gap-4 md:h-24 md:flex-row">
          <div className="flex flex-col items-center gap-4 px-8 md:flex-row md:gap-2 md:px-0">
            <p className="text-center text-sm leading-loose text-muted-foreground md:text-left">
              &copy; {new Date().getFullYear()} GoFrame. All rights reserved.
            </p>
          </div>
          <div className="flex items-center">
            <Link
              href={siteConfig.links.github}
              target="_blank"
              rel="noreferrer"
              className="px-4 py-2 hover:text-primary"
            >
              GitHub
            </Link>
            <Link
              href={siteConfig.links.twitter}
              target="_blank"
              rel="noreferrer"
              className="px-4 py-2 hover:text-primary"
            >
              Twitter
            </Link>
          </div>
        </div>
      </footer>
    </div>
  )
}

