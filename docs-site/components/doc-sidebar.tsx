"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"

import { cn } from "@/lib/utils"
import { ScrollArea } from "@/components/ui/scroll-area"

interface DocSidebarProps {
  items: {
    title: string
    href: string
    items?: {
      title: string
      href: string
    }[]
  }[]
}

export function DocSidebar({ items }: DocSidebarProps) {
  const pathname = usePathname()

  return (
    <ScrollArea className="h-[calc(100vh-3.5rem)] py-6 pr-6 lg:py-8">
      <div className="w-full">
        {items.map((item) => (
          <div key={item.href} className="pb-4">
            <h4 className="mb-1 rounded-md px-2 py-1 text-sm font-semibold">{item.title}</h4>
            {item.items?.length && (
              <div className="grid grid-flow-row auto-rows-max text-sm">
                {item.items.map((subItem) => (
                  <Link
                    key={subItem.href}
                    href={subItem.href}
                    className={cn(
                      "group flex w-full items-center rounded-md border border-transparent px-2 py-1 hover:underline",
                      pathname === subItem.href ? "font-medium text-primary" : "text-muted-foreground",
                    )}
                  >
                    {subItem.title}
                  </Link>
                ))}
              </div>
            )}
          </div>
        ))}
      </div>
    </ScrollArea>
  )
}

