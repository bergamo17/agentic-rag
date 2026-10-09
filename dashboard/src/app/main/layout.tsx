"use client";

import { Suspense, useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { ConversationList } from "@/components/chat/conversations-list";
import { MessageSquare, LayoutGrid, ChevronLeft, ChevronRight } from "lucide-react";

const navItems = [
  { href: "/main/chat", label: "Chat", icon: MessageSquare },
  { href: "/main/dashboard", label: "Dashboard", icon: LayoutGrid },
];

export default function MainLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const [collapsed, setCollapsed] = useState(false);

  return (
    <div className="flex h-screen">
      <aside
        className={`relative shrink-0 flex flex-col gap-1 p-4 bg-sidebar text-grey-50 transition-all duration-300 ease-in-out ${
          collapsed ? "w-16" : "w-56"
        }`}
      >
        <button
          onClick={() => setCollapsed((prev) => !prev)}
          className="absolute -right-3 top-6 z-10 flex h-6 w-6 items-center justify-center
                     rounded-full border border-grey-200 bg-white shadow-sm
                     hover:bg-grey-50 transition-colors"
          aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
        >
          {collapsed ? (
            <ChevronRight size={14} className="text-grey-600" />
          ) : (
            <ChevronLeft size={14} className="text-grey-600" />
          )}
        </button>

        <div
          className={`mb-4 px-3 text-sm font-semibold text-white whitespace-nowrap overflow-hidden transition-opacity duration-200 ${
            collapsed ? "opacity-0" : "opacity-100"
          }`}
        >
          Agentic RAG
        </div>

        {navItems.map(({ href, label, icon: Icon }) => {
          const isActive = pathname === href;
          return (
            <Link
              key={href}
              href={href}
              title={collapsed ? label : undefined}
              className={`relative flex items-center gap-2 rounded-md px-3 py-2 text-sm font-medium transition-colors ${
                isActive
                  ? "bg-white/15 text-white"
                  : "text-[#C7E1F7] hover:bg-white/10"
              }`}
            >
              {isActive && (
                <span className="absolute left-0 top-1/2 h-5 w-[3px] -translate-y-1/2 rounded-r bg-accent shadow-[0_0_8px_rgba(14,213,221,0.55)]" />
              )}
              <Icon size={16} className="shrink-0" />
              {!collapsed && <span className="whitespace-nowrap">{label}</span>}
            </Link>
          );
        })}

        {!collapsed && pathname.startsWith("/main/chat") && (
          <div className="mt-4 min-h-0 flex-1 overflow-y-auto">
            <Suspense fallback={null}>
              <ConversationList />
            </Suspense>
          </div>
        )}
      </aside>
      <main className="flex-1 overflow-hidden bg-background">{children}</main>
    </div>
  );
}