import Link from "next/link";

export default function MainLayout({ children }: { children: React.ReactNode }) {
    return (
      <div className="flex h-screen flex-col">
        <nav className="flex gap-4 border-b p-3">
            <Link href="/chat" className="text-sm font-medium">Chat</Link>
            <Link href="/dashboard" className="text-sm font-medium">Dashboard</Link>
        </nav>
        <main className="flex-1 overflow-hidden">{children}</main>
      </div>
    );
}