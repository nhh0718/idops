"use client";

import {
  Activity,
  Clock,
  Container,
  Cpu,
  FileCode,
  HardDrive,
  KeyRound,
  Search,
  Server,
  Shield,
  TrendingUp,
  Wifi,
  Zap,
} from "lucide-react";
import { useEffect, useState } from "react";
import type { DockerContainer, PortEntry, SSHHost } from "../types";
import type { TabId } from "./Sidebar";

interface OverviewTabProps {
  containers: DockerContainer[];
  ports: PortEntry[];
  sshHosts: SSHHost[];
  envVarCount: number;
  onNavigate: (tab: TabId) => void;
}

export default function OverviewTab({
  containers,
  ports,
  sshHosts,
  envVarCount,
  onNavigate,
}: OverviewTabProps) {
  const [now, setNow] = useState("");
  useEffect(() => {
    setNow(new Date().toLocaleString("vi-VN"));
    const interval = setInterval(
      () => setNow(new Date().toLocaleString("vi-VN")),
      1000,
    );
    return () => clearInterval(interval);
  }, []);

  const runningContainers = containers.filter(
    (c) => c.state === "running",
  ).length;
  const stoppedContainers = containers.length - runningContainers;
  const avgCpu =
    containers
      .filter((c) => c.cpu !== undefined)
      .reduce((sum, c) => sum + (c.cpu || 0), 0) /
    Math.max(containers.filter((c) => c.cpu !== undefined).length, 1);
  const avgMem =
    containers
      .filter((c) => c.memory !== undefined)
      .reduce((sum, c) => sum + (c.memory || 0), 0) /
    Math.max(containers.filter((c) => c.memory !== undefined).length, 1);
  const tcpPorts = ports.filter((p) => p.protocol === "tcp").length;
  const connectedHosts = sshHosts.filter(
    (h) => h.status === "connected",
  ).length;

  const statCards = [
    {
      title: "Docker Containers",
      value: containers.length.toString(),
      sub: `${runningContainers} running, ${stoppedContainers} stopped`,
      icon: <Container size={22} />,
      color: "from-purple-500/20 to-purple-600/10",
      iconBg: "bg-purple-500/20 text-purple-400",
      tab: "docker" as TabId,
    },
    {
      title: "Active Ports",
      value: ports.length.toString(),
      sub: `${tcpPorts} TCP, ${ports.length - tcpPorts} UDP`,
      icon: <Search size={22} />,
      color: "from-blue-500/20 to-blue-600/10",
      iconBg: "bg-blue-500/20 text-blue-400",
      tab: "ports" as TabId,
    },
    {
      title: "SSH Hosts",
      value: sshHosts.length.toString(),
      sub: `${connectedHosts} reachable`,
      icon: <KeyRound size={22} />,
      color: "from-emerald-500/20 to-emerald-600/10",
      iconBg: "bg-emerald-500/20 text-emerald-400",
      tab: "ssh" as TabId,
    },
    {
      title: "Env Variables",
      value: envVarCount.toString(),
      sub: "Across all .env files",
      icon: <FileCode size={22} />,
      color: "from-amber-500/20 to-amber-600/10",
      iconBg: "bg-amber-500/20 text-amber-400",
      tab: "env" as TabId,
    },
  ];

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-white">Dashboard Overview</h2>
          <p className="text-sm text-[var(--color-muted)] mt-1">
            System status at a glance — {now}
          </p>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-emerald-500/10 rounded-full border border-emerald-500/20">
          <div className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse-slow" />
          <span className="text-xs text-emerald-400 font-medium">
            System Online
          </span>
        </div>
      </div>

      {/* Stat Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4">
        {statCards.map((card) => (
          <button
            key={card.title}
            onClick={() => onNavigate(card.tab)}
            className={`bg-gradient-to-br ${card.color} border border-[var(--color-border)] rounded-xl p-5 text-left hover:border-[var(--color-primary)] transition-all duration-200 group`}
          >
            <div className="flex items-start justify-between">
              <div>
                <p className="text-xs text-[var(--color-muted)] font-medium uppercase tracking-wider">
                  {card.title}
                </p>
                <p className="text-3xl font-bold text-white mt-1">
                  {card.value}
                </p>
                <p className="text-xs text-[var(--color-muted)] mt-1">
                  {card.sub}
                </p>
              </div>
              <div className={`p-2.5 rounded-lg ${card.iconBg}`}>
                {card.icon}
              </div>
            </div>
          </button>
        ))}
      </div>

      {/* Resource Usage + Quick Actions */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        {/* Resource Gauges */}
        <div className="lg:col-span-2 bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl p-5">
          <h3 className="text-sm font-semibold text-white mb-4 flex items-center gap-2">
            <Activity size={16} className="text-[var(--color-primary)]" />
            Resource Usage
          </h3>
          <div className="grid grid-cols-2 gap-6">
            <ResourceGauge
              label="Avg CPU"
              value={avgCpu}
              icon={<Cpu size={16} />}
              color="purple"
            />
            <ResourceGauge
              label="Avg Memory"
              value={avgMem}
              icon={<HardDrive size={16} />}
              color="blue"
            />
            <ResourceGauge
              label="Network Health"
              value={connectedHosts > 0 ? 95 : 0}
              icon={<Wifi size={16} />}
              color="emerald"
            />
            <ResourceGauge
              label="System Load"
              value={Math.min((ports.length / 50) * 100, 100)}
              icon={<TrendingUp size={16} />}
              color="amber"
            />
          </div>
        </div>

        {/* Quick Actions */}
        <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl p-5">
          <h3 className="text-sm font-semibold text-white mb-4 flex items-center gap-2">
            <Zap size={16} className="text-amber-400" />
            Quick Actions
          </h3>
          <div className="space-y-2">
            {[
              {
                label: "Scan Ports",
                icon: <Search size={14} />,
                tab: "ports" as TabId,
                color: "text-blue-400",
              },
              {
                label: "Docker Dashboard",
                icon: <Container size={14} />,
                tab: "docker" as TabId,
                color: "text-purple-400",
              },
              {
                label: "Manage SSH Hosts",
                icon: <KeyRound size={14} />,
                tab: "ssh" as TabId,
                color: "text-emerald-400",
              },
              {
                label: "Env File Manager",
                icon: <FileCode size={14} />,
                tab: "env" as TabId,
                color: "text-amber-400",
              },
              {
                label: "Generate Nginx Config",
                icon: <Server size={14} />,
                tab: "nginx" as TabId,
                color: "text-rose-400",
              },
            ].map((action) => (
              <button
                key={action.label}
                onClick={() => onNavigate(action.tab)}
                className="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg bg-[var(--color-background)] hover:bg-[var(--color-card-hover)] border border-transparent hover:border-[var(--color-border)] transition-all text-left group"
              >
                <span className={action.color}>{action.icon}</span>
                <span className="text-sm text-[var(--color-muted)] group-hover:text-white transition-colors">
                  {action.label}
                </span>
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Recent Activity */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        {/* Recent Containers */}
        <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl p-5">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-sm font-semibold text-white flex items-center gap-2">
              <Container size={16} className="text-purple-400" />
              Recent Containers
            </h3>
            <button
              onClick={() => onNavigate("docker")}
              className="text-xs text-[var(--color-primary)] hover:text-[var(--color-primary-light)]"
            >
              View All
            </button>
          </div>
          {containers.length === 0 ? (
            <p className="text-sm text-[var(--color-muted)] text-center py-4">
              No containers found
            </p>
          ) : (
            <div className="space-y-2">
              {containers.slice(0, 5).map((c) => (
                <div
                  key={c.id}
                  className="flex items-center justify-between py-2 px-3 rounded-lg bg-[var(--color-background)] border border-[var(--color-border)]"
                >
                  <div className="flex items-center gap-3">
                    <div
                      className={`w-2 h-2 rounded-full ${c.state === "running" ? "bg-emerald-400" : "bg-red-400"}`}
                    />
                    <div>
                      <p className="text-sm text-white font-medium">{c.name}</p>
                      <p className="text-[10px] text-[var(--color-muted)]">
                        {c.image}
                      </p>
                    </div>
                  </div>
                  <span
                    className={`text-xs px-2 py-0.5 rounded-full ${c.state === "running" ? "bg-emerald-500/10 text-emerald-400" : "bg-red-500/10 text-red-400"}`}
                  >
                    {c.state}
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Active Ports */}
        <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl p-5">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-sm font-semibold text-white flex items-center gap-2">
              <Search size={16} className="text-blue-400" />
              Active Ports
            </h3>
            <button
              onClick={() => onNavigate("ports")}
              className="text-xs text-[var(--color-primary)] hover:text-[var(--color-primary-light)]"
            >
              View All
            </button>
          </div>
          {ports.length === 0 ? (
            <p className="text-sm text-[var(--color-muted)] text-center py-4">
              No active ports
            </p>
          ) : (
            <div className="space-y-2">
              {ports.slice(0, 5).map((p, i) => (
                <div
                  key={i}
                  className="flex items-center justify-between py-2 px-3 rounded-lg bg-[var(--color-background)] border border-[var(--color-border)]"
                >
                  <div className="flex items-center gap-3">
                    <span
                      className={`text-xs px-1.5 py-0.5 rounded font-mono ${p.protocol === "tcp" ? "bg-blue-500/10 text-blue-400" : "bg-amber-500/10 text-amber-400"}`}
                    >
                      {p.protocol.toUpperCase()}
                    </span>
                    <div>
                      <p className="text-sm text-white font-mono">:{p.port}</p>
                      <p className="text-[10px] text-[var(--color-muted)]">
                        {p.process}
                      </p>
                    </div>
                  </div>
                  <span className="text-xs text-[var(--color-muted)] font-mono">
                    PID {p.pid}
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Footer Info */}
      <div className="flex items-center justify-between py-3 px-4 bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl">
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2 text-xs text-[var(--color-muted)]">
            <Clock size={12} />
            <span>Last refresh: {now}</span>
          </div>
          <div className="flex items-center gap-2 text-xs text-[var(--color-muted)]">
            <Shield size={12} />
            <span>idops v1.0.0</span>
          </div>
        </div>
        <div className="flex items-center gap-2 text-xs text-emerald-400">
          <div className="w-1.5 h-1.5 rounded-full bg-emerald-400" />
          All systems operational
        </div>
      </div>
    </div>
  );
}

function ResourceGauge({
  label,
  value,
  icon,
  color,
}: {
  label: string;
  value: number;
  icon: React.ReactNode;
  color: string;
}) {
  const colorMap: Record<string, { bar: string; text: string }> = {
    purple: { bar: "bg-purple-500", text: "text-purple-400" },
    blue: { bar: "bg-blue-500", text: "text-blue-400" },
    emerald: { bar: "bg-emerald-500", text: "text-emerald-400" },
    amber: { bar: "bg-amber-500", text: "text-amber-400" },
  };
  const c = colorMap[color] || colorMap.purple;
  const clamped = Math.min(Math.max(value, 0), 100);

  return (
    <div>
      <div className="flex items-center justify-between mb-2">
        <div className="flex items-center gap-2">
          <span className={c.text}>{icon}</span>
          <span className="text-xs text-[var(--color-muted)]">{label}</span>
        </div>
        <span className={`text-sm font-bold ${c.text}`}>
          {clamped.toFixed(1)}%
        </span>
      </div>
      <div className="h-2 bg-[var(--color-background)] rounded-full overflow-hidden">
        <div
          className={`h-full ${c.bar} rounded-full transition-all duration-500`}
          style={{ width: `${clamped}%` }}
        />
      </div>
    </div>
  );
}
