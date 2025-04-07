import Link from "next/link"
import { ChevronLeft, ChevronRight } from "lucide-react"

import { cn } from "@/lib/utils"
import { buttonVariants } from "@/components/ui/button"

interface DocPaginationProps {
  prev?: {
    title: string
    href: string
  }
  next?: {
    title: string
    href: string
  }
}

export function DocPagination({ prev, next }: DocPaginationProps) {
  return (
    <div className="flex flex-row items-center justify-between">
      {prev && (
        <Link href={prev.href} className={cn(buttonVariants({ variant: "ghost" }), "gap-1")}>
          <ChevronLeft className="h-4 w-4" />
          {prev.title}
        </Link>
      )}
      {next && (
        <Link href={next.href} className={cn(buttonVariants({ variant: "ghost" }), "ml-auto gap-1")}>
          {next.title}
          <ChevronRight className="h-4 w-4" />
        </Link>
      )}
    </div>
  )
}

