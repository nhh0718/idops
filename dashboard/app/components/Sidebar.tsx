"use client";

import {
  LayoutDashboard,
  Container,
  Search,
  KeyRound,
  FileCode,
  Server,
  Terminal,
  ChevronLeft,
  ChevronRight,
} from "lucide-react";

export type TabId = "overview" | "docker" | "ports" | "ssh" | "env" | "nginx";

interface SidebarProps {
  activeTab: TabId;
  onTabChange: (tab: TabId) => void;
  collapsed: boolean;
  onToggle: () => void;
}

const menuItems: { id: TabId; label: string; icon: React.ReactNode; desc: string }[] = [
  { id: "overview", label: "Overview", icon: <LayoutDashboard size={20} />, desc: "System overview" },
  { id: "docker", label: "Docker", icon: <Container size={20} />, desc: "Container management" },
  { id: "ports", label: "Ports", icon: <Search size={20} />, desc: "Port scanner" },
  { id: "ssh", label: "SSH", icon: <KeyRound size={20} />, desc: "SSH host manager" },
  { id: "env", label: "Env", icon: <FileCode size={20} />, desc: "Environment files" },
  { id: "nginx", label: "Nginx", icon: <Server size={20} />, desc: "Config generator" },
];

export default function Sidebar({ activeTab, onTabChange, collapsed, onToggle }: SidebarProps) {
  return (
    <aside
      className={`fixed left-0 top-0 h-full bg-[var(--color-sidebar)] border-r border-[var(--color-border)] z-50 flex flex-col transition-all duration-300 ${
        collapsed ? "w-[68px]" : "w-[240px]"
      }`}
    >
      <div className="flex items-center gap-3 px-4 h-16 border-b border-[var(--color-border)]">
        <div className="w-8 h-8 rounded-lg bg-[var(--color-primary)] flex items-center justify-center flex-shrink-0">
          <Terminal size={18} className="text-white" />
        </div>
        {!collapsed && (
          <div className="animate-fade-in">
            <h1 className="text-base font-bold text-white tracking-tight">idops</h1>
            <p className="text-[10px] text-[var(--color-muted)] -mt-0.5">DevOps Toolkit</p>
          </div>
        )}
      </div>

      <nav className="flex-1 py-3 px-2 space-y-1 overflow-y-auto">
        {menuItems.map((item) => (
          <button
            key={item.id}
            onClick={() => onTabChange(item.id)}
            className={`w-full flex items-center gap-3 px-3 py-2.5 rounded-lg transition-all duration-200 group ${
              activeTab === item.id
                ? "bg-[var(--color-primary)] text-white shadow-lg shadow-purple-500/20"
                : "text-[var(--color-muted)] hover:bg-[var(--color-card)] hover:text-white"
            }`}
            title={collapsed ? item.label : undefined}
          >
            <span className="flex-shrink-0">{item.icon}</span>
            {!collapsed && (
              <div className="text-left animate-fade-in">
                <div className="text-sm font-medium">{item.label}</div>
                <div className={`text-[10px] ${activeTab === item.id ? "text-purple-200" : "text-[var(--color-muted)]"}`}>
                  {item.desc}
                </div>
              </div>
            )}
          </button>
        ))}
      </nav>

      <button
        onClick={onToggle}
        className="flex items-center justify-center h-12 border-t border-[var(--color-border)] text-[var(--color-muted)] hover:text-white transition-colors"
      >
        {collapsed ? <ChevronRight size={18} /> : <ChevronLeft size={18} />}
      </button>
    </aside>
  );
}
