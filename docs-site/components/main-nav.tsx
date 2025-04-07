"use client"

import * as React from "react"
import Link from "next/link"
import { usePathname } from "next/navigation"

import { cn } from "@/lib/utils"
import { Icons } from "@/components/icons"
import { MobileNav } from "@/components/mobile-nav"

export function MainNav() {
  const pathname = usePathname()
  const [showMobileMenu, setShowMobileMenu] = React.useState<boolean>(false)

  return (
    <div className="flex md:gap-10 gap-6 items-center">
      <Link href="/" className="hidden md:block">
        <span className="font-bold text-xl flex items-center gap-2">
          <Icons.logo className="h-6 w-6" />
          <span>GoFrame</span>
        </span>
        <span className="sr-only">GoFrame</span>
      </Link>
      <nav className="hidden md:flex gap-6">
        <Link
          href="/docs"
          className={cn(
            "text-sm font-medium transition-colors hover:text-primary",
            pathname === "/docs" || pathname?.startsWith("/docs/") ? "text-primary" : "text-muted-foreground",
          )}
        >
          Documentation
        </Link>
        <Link
          href="/docs/api"
          className={cn(
            "text-sm font-medium transition-colors hover:text-primary",
            pathname === "/docs/api" ? "text-primary" : "text-muted-foreground",
          )}
        >
          API
        </Link>
        <Link
          href="/docs/examples"
          className={cn(
            "text-sm font-medium transition-colors hover:text-primary",
            pathname === "/docs/examples" ? "text-primary" : "text-muted-foreground",
          )}
        >
          Examples
        </Link>
        <Link
          href="/docs/community"
          className={cn(
            "text-sm font-medium transition-colors hover:text-primary",
            pathname === "/docs/community" ? "text-primary" : "text-muted-foreground",
          )}
        >
          Community
        </Link>
      </nav>
      <button className="flex items-center space-x-2 md:hidden" onClick={() => setShowMobileMenu(!showMobileMenu)}>
        {showMobileMenu ? <Icons.close className="h-5 w-5" /> : <Icons.menu className="h-5 w-5" />}
        <span className="font-bold">Menu</span>
      </button>
      {showMobileMenu && <MobileNav />}
    </div>
  )
}

