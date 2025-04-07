"use client"

import * as React from "react"
import { useRouter } from "next/navigation"
import type { DialogProps } from "@radix-ui/react-dialog"
import { SearchIcon } from "lucide-react"

import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command"

interface SearchProps extends DialogProps {}

export function Search({ ...props }: SearchProps) {
  const router = useRouter()
  const [open, setOpen] = React.useState(false)
  const [query, setQuery] = React.useState("")
  const [results, setResults] = React.useState<SearchResult[]>([])

  React.useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault()
        setOpen((open) => !open)
      }
    }

    document.addEventListener("keydown", down)
    return () => document.removeEventListener("keydown", down)
  }, [])

  const handleSearch = React.useCallback(async (value: string) => {
    setQuery(value)
    if (value.length === 0) {
      setResults([])
      return
    }

    // Perform search
    const searchResults = await performSearch(value)
    setResults(searchResults)
  }, [])

  const handleSelect = React.useCallback((callback: () => unknown) => {
    setOpen(false)
    callback()
  }, [])

  return (
    <>
      <Button
        variant="outline"
        className={cn(
          "relative h-9 w-full justify-start rounded-[0.5rem] text-sm text-muted-foreground sm:pr-12 md:w-40 lg:w-64",
        )}
        onClick={() => setOpen(true)}
        {...props}
      >
        <SearchIcon className="mr-2 h-4 w-4" />
        <span className="hidden lg:inline-flex">Search documentation...</span>
        <span className="inline-flex lg:hidden">Search...</span>
        <kbd className="pointer-events-none absolute right-1.5 top-1.5 hidden h-6 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100 sm:flex">
          <span className="text-xs">⌘</span>K
        </kbd>
      </Button>
      <CommandDialog open={open} onOpenChange={setOpen}>
        <CommandInput placeholder="Search documentation..." value={query} onValueChange={handleSearch} />
        <CommandList>
          <CommandEmpty>No results found.</CommandEmpty>
          <CommandGroup heading="Documentation">
            {results.map((result) => (
              <CommandItem key={result.href} onSelect={() => handleSelect(() => router.push(result.href))}>
                <span>{result.title}</span>
              </CommandItem>
            ))}
          </CommandGroup>
        </CommandList>
      </CommandDialog>
    </>
  )
}

interface SearchResult {
  title: string
  href: string
  content?: string
}

// This would typically fetch from an API or search index
// For simplicity, we're implementing a client-side search
async function performSearch(query: string): Promise<SearchResult[]> {
  // This is a simplified search implementation
  // In a real app, you'd use a proper search index or API
  const searchIndex: SearchResult[] = [
    { title: "Introduction", href: "/docs" },
    { title: "Installation", href: "/docs/installation" },
    { title: "Project Structure", href: "/docs/project-structure" },
    { title: "Configuration", href: "/docs/configuration" },
    { title: "Routing", href: "/docs/core-concepts/routing" },
    { title: "Controllers", href: "/docs/core-concepts/controllers" },
    { title: "Models", href: "/docs/core-concepts/models" },
    { title: "Middleware", href: "/docs/core-concepts/middleware" },
    { title: "Database", href: "/docs/database" },
    { title: "Migrations", href: "/docs/database/migrations" },
    { title: "Query Builder", href: "/docs/database/query-builder" },
    { title: "Repositories", href: "/docs/database/repositories" },
    { title: "Authentication", href: "/docs/features/authentication" },
    { title: "Resources", href: "/docs/features/resources" },
    { title: "Rate Limiting", href: "/docs/features/rate-limiting" },
    { title: "CLI Commands", href: "/docs/features/cli-commands" },
  ]

  // Filter results based on query
  return searchIndex.filter((item) => item.title.toLowerCase().includes(query.toLowerCase()))
}

