"use client"
import Link from "next/link"

import { cn } from "@/lib/utils"
import { Icons } from "@/components/icons"
import { useLockBody } from "@/hooks/use-lock-body"

export function MobileNav() {
  useLockBody()

  return (
    <div
      className={cn(
        "fixed inset-0 top-16 z-50 grid h-[calc(100vh-4rem)] grid-flow-row auto-rows-max overflow-auto p-6 pb-32 shadow-md animate-in slide-in-from-bottom-80 md:hidden",
      )}
    >
      <div className="relative z-20 grid gap-6 rounded-md bg-card p-4 text-popover-foreground shadow-md">
        <Link href="/" className="flex items-center space-x-2">
          <Icons.logo className="h-6 w-6" />
          <span className="font-bold">GoFrame</span>
        </Link>
        <nav className="grid grid-flow-row auto-rows-max text-sm">
          <Link href="/docs" className="flex w-full items-center rounded-md p-2 text-sm font-medium hover:underline">
            Documentation
          </Link>
          <Link
            href="/docs/api"
            className="flex w-full items-center rounded-md p-2 text-sm font-medium hover:underline"
          >
            API
          </Link>
          <Link
            href="/docs/examples"
            className="flex w-full items-center rounded-md p-2 text-sm font-medium hover:underline"
          >
            Examples
          </Link>
          <Link
            href="/docs/community"
            className="flex w-full items-center rounded-md p-2 text-sm font-medium hover:underline"
          >
            Community
          </Link>
        </nav>
      </div>
    </div>
  )
}

