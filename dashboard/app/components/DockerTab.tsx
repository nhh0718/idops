"use client";

import {
  ArrowDownToLine,
  ArrowUpFromLine,
  Container,
  Cpu,
  HardDrive,
  Info,
  Play,
  RefreshCw,
  RotateCw,
  ScrollText,
  Search,
  Square,
  Trash2,
  X,
} from "lucide-react";
import { useEffect, useState } from "react";
import { dockerApi } from "../lib/api";
import type { DockerContainer } from "../types";

interface DockerTabProps {
  containers: DockerContainer[];
  setContainers: React.Dispatch<React.SetStateAction<DockerContainer[]>>;
}

export default function DockerTab({
  containers: initialContainers,
}: {
  containers: DockerContainer[];
}) {
  const [containers, setContainers] =
    useState<DockerContainer[]>(initialContainers);
  const [filter, setFilter] = useState("");
  const [selectedContainer, setSelectedContainer] =
    useState<DockerContainer | null>(null);
  const [showLogs, setShowLogs] = useState(false);
  const [showInspect, setShowInspect] = useState(false);
  const [confirmAction, setConfirmAction] = useState<{
    id: string;
    action: string;
  } | null>(null);
  const [statusMsg, setStatusMsg] = useState<{
    text: string;
    isError: boolean;
  } | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [logsContent, setLogsContent] = useState("");

  useEffect(() => {
    loadContainers();
  }, []);

  async function loadContainers() {
    setIsLoading(true);
    try {
      const data = await dockerApi.list();
      setContainers(data);
    } catch (error) {
      showStatus("Failed to load containers", true);
    } finally {
      setIsLoading(false);
    }
  }

  async function handleRefresh() {
    await loadContainers();
    showStatus("Refreshed container list");
  }

  const filtered = containers.filter(
    (c) =>
      c.name.toLowerCase().includes(filter.toLowerCase()) ||
      c.image.toLowerCase().includes(filter.toLowerCase()) ||
      c.id.toLowerCase().includes(filter.toLowerCase()),
  );

  const running = containers.filter((c) => c.state === "running").length;
  const stopped = containers.filter((c) => c.state === "exited").length;

  function showStatus(text: string, isError = false) {
    setStatusMsg({ text, isError });
    setTimeout(() => setStatusMsg(null), 3000);
  }

  function handleAction(id: string, action: string) {
    const container = containers.find((c) => c.id === id);
    if (!container) return;

    if (action === "remove" || action === "stop") {
      setConfirmAction({ id, action });
      return;
    }
    executeAction(id, action);
  }

  async function executeAction(id: string, action: string) {
    setConfirmAction(null);
    const container = containers.find((c) => c.id === id);
    if (!container) return;

    try {
      const result = await dockerApi.action(action, id);
      if (result.error) {
        showStatus(result.error, true);
      } else {
        showStatus(result.message || `${action} completed`);
        await loadContainers();
      }
    } catch (error) {
      showStatus(`Failed to ${action} container`, true);
    }
  }

  async function handleShowLogs(container: DockerContainer) {
    setSelectedContainer(container);
    setShowLogs(true);
    setLogsContent("Loading logs...");
    try {
      const logs = await dockerApi.logs(container.id);
      setLogsContent(logs || "No logs available");
    } catch {
      setLogsContent("Failed to load logs");
    }
  }

  return (
    <div className="space-y-4 animate-fade-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-white flex items-center gap-3">
            <Container size={24} className="text-purple-400" />
            Docker Dashboard
          </h2>
          <p className="text-sm text-[var(--color-muted)] mt-1">
            Manage containers, monitor resources, view logs
          </p>
        </div>
        <button
          onClick={handleRefresh}
          disabled={isLoading}
          className="flex items-center gap-2 px-3 py-2 bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-muted)] hover:text-white hover:border-[var(--color-primary)] transition-all disabled:opacity-50"
        >
          <RefreshCw size={14} className={isLoading ? "animate-spin" : ""} />
          Refresh
        </button>
      </div>

      {/* Stats Bar */}
      <div className="flex gap-4">
        <div className="flex items-center gap-2 px-3 py-1.5 bg-[var(--color-card)] rounded-lg border border-[var(--color-border)]">
          <div className="w-2 h-2 rounded-full bg-emerald-400" />
          <span className="text-xs text-[var(--color-muted)]">Running:</span>
          <span className="text-xs font-bold text-emerald-400">{running}</span>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-[var(--color-card)] rounded-lg border border-[var(--color-border)]">
          <div className="w-2 h-2 rounded-full bg-red-400" />
          <span className="text-xs text-[var(--color-muted)]">Stopped:</span>
          <span className="text-xs font-bold text-red-400">{stopped}</span>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-[var(--color-card)] rounded-lg border border-[var(--color-border)]">
          <Container size={12} className="text-[var(--color-muted)]" />
          <span className="text-xs text-[var(--color-muted)]">Total:</span>
          <span className="text-xs font-bold text-white">
            {containers.length}
          </span>
        </div>
      </div>

      {/* Status Message */}
      {statusMsg && (
        <div
          className={`px-4 py-2 rounded-lg text-sm ${statusMsg.isError ? "bg-red-500/10 text-red-400 border border-red-500/20" : "bg-emerald-500/10 text-emerald-400 border border-emerald-500/20"}`}
        >
          {statusMsg.text}
        </div>
      )}

      {/* Filter */}
      <div className="relative">
        <Search
          size={16}
          className="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--color-muted)]"
        />
        <input
          type="text"
          placeholder="Filter by name, image, or ID..."
          value={filter}
          onChange={(e) => setFilter(e.target.value)}
          className="w-full pl-10 pr-4 py-2.5 bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg text-sm text-white placeholder:text-[var(--color-muted)] focus:outline-none focus:border-[var(--color-primary)] transition-colors"
        />
      </div>

      {/* Table */}
      <div className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-[var(--color-border)]">
                <th className="px-4 py-3 text-left text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider">
                  Container
                </th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider">
                  Image
                </th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider">
                  State
                </th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider">
                  Status
                </th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider">
                  CPU%
                </th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider">
                  Mem%
                </th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider">
                  Net I/O
                </th>
                <th className="px-4 py-3 text-right text-xs font-semibold text-[var(--color-muted)] uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody>
              {filtered.map((c) => (
                <tr
                  key={c.id}
                  className="border-b border-[var(--color-border)] hover:bg-[var(--color-card-hover)] transition-colors"
                >
                  <td className="px-4 py-3">
                    <div>
                      <p className="font-medium text-white">{c.name}</p>
                      <p className="text-[10px] text-[var(--color-muted)] font-mono">
                        {c.id}
                      </p>
                    </div>
                  </td>
                  <td className="px-4 py-3 text-[var(--color-muted)] font-mono text-xs">
                    {c.image}
                  </td>
                  <td className="px-4 py-3">
                    <span
                      className={`inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium ${
                        c.state === "running"
                          ? "bg-emerald-500/10 text-emerald-400"
                          : c.state === "paused"
                            ? "bg-amber-500/10 text-amber-400"
                            : "bg-red-500/10 text-red-400"
                      }`}
                    >
                      <div
                        className={`w-1.5 h-1.5 rounded-full ${
                          c.state === "running"
                            ? "bg-emerald-400"
                            : c.state === "paused"
                              ? "bg-amber-400"
                              : "bg-red-400"
                        }`}
                      />
                      {c.state}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-xs text-[var(--color-muted)]">
                    {c.status}
                  </td>
                  <td className="px-4 py-3">
                    {c.cpu !== undefined ? (
                      <div className="flex items-center gap-2">
                        <div className="w-16 h-1.5 bg-[var(--color-background)] rounded-full overflow-hidden">
                          <div
                            className={`h-full rounded-full ${c.cpu > 80 ? "bg-red-500" : c.cpu > 50 ? "bg-amber-500" : "bg-emerald-500"}`}
                            style={{ width: `${Math.min(c.cpu, 100)}%` }}
                          />
                        </div>
                        <span className="text-xs text-[var(--color-muted)]">
                          {c.cpu.toFixed(1)}%
                        </span>
                      </div>
                    ) : (
                      <span className="text-xs text-[var(--color-muted)]">
                        -
                      </span>
                    )}
                  </td>
                  <td className="px-4 py-3">
                    {c.memory !== undefined ? (
                      <div className="flex items-center gap-2">
                        <div className="w-16 h-1.5 bg-[var(--color-background)] rounded-full overflow-hidden">
                          <div
                            className={`h-full rounded-full ${c.memory > 80 ? "bg-red-500" : c.memory > 50 ? "bg-amber-500" : "bg-blue-500"}`}
                            style={{ width: `${Math.min(c.memory, 100)}%` }}
                          />
                        </div>
                        <span className="text-xs text-[var(--color-muted)]">
                          {c.memory.toFixed(1)}%
                        </span>
                      </div>
                    ) : (
                      <span className="text-xs text-[var(--color-muted)]">
                        -
                      </span>
                    )}
                  </td>
                  <td className="px-4 py-3 text-xs text-[var(--color-muted)]">
                    {c.netIn && c.netOut ? (
                      <div className="flex items-center gap-2">
                        <ArrowDownToLine
                          size={10}
                          className="text-emerald-400"
                        />
                        <span>{c.netIn}</span>
                        <ArrowUpFromLine size={10} className="text-blue-400" />
                        <span>{c.netOut}</span>
                      </div>
                    ) : (
                      "-"
                    )}
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex items-center justify-end gap-1">
                      {c.state !== "running" && (
                        <button
                          onClick={() => handleAction(c.id, "start")}
                          title="Start"
                          className="p-1.5 rounded hover:bg-emerald-500/10 text-[var(--color-muted)] hover:text-emerald-400 transition-colors"
                        >
                          <Play size={14} />
                        </button>
                      )}
                      {c.state === "running" && (
                        <button
                          onClick={() => handleAction(c.id, "stop")}
                          title="Stop"
                          className="p-1.5 rounded hover:bg-red-500/10 text-[var(--color-muted)] hover:text-red-400 transition-colors"
                        >
                          <Square size={14} />
                        </button>
                      )}
                      <button
                        onClick={() => handleAction(c.id, "restart")}
                        title="Restart"
                        className="p-1.5 rounded hover:bg-amber-500/10 text-[var(--color-muted)] hover:text-amber-400 transition-colors"
                      >
                        <RotateCw size={14} />
                      </button>
                      <button
                        onClick={() => handleShowLogs(c)}
                        title="Logs"
                        className="p-1.5 rounded hover:bg-blue-500/10 text-[var(--color-muted)] hover:text-blue-400 transition-colors"
                      >
                        <ScrollText size={14} />
                      </button>
                      <button
                        onClick={() => {
                          setSelectedContainer(c);
                          setShowInspect(true);
                        }}
                        title="Inspect"
                        className="p-1.5 rounded hover:bg-purple-500/10 text-[var(--color-muted)] hover:text-purple-400 transition-colors"
                      >
                        <Info size={14} />
                      </button>
                      <button
                        onClick={() => handleAction(c.id, "remove")}
                        title="Remove"
                        className="p-1.5 rounded hover:bg-red-500/10 text-[var(--color-muted)] hover:text-red-400 transition-colors"
                      >
                        <Trash2 size={14} />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
              {filtered.length === 0 && (
                <tr>
                  <td
                    colSpan={8}
                    className="px-4 py-8 text-center text-[var(--color-muted)]"
                  >
                    No containers found
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Keyboard shortcuts info */}
      <div className="flex flex-wrap gap-3 text-xs text-[var(--color-muted)]">
        <span className="flex items-center gap-1">
          <Cpu size={10} /> CPU/Memory stats auto-refresh
        </span>
        <span>•</span>
        <span className="flex items-center gap-1">
          <HardDrive size={10} /> Click container for details
        </span>
      </div>

      {/* Confirm Dialog */}
      {confirmAction && (
        <div
          className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
          onClick={() => setConfirmAction(null)}
        >
          <div
            className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl p-6 max-w-sm w-full mx-4"
            onClick={(e) => e.stopPropagation()}
          >
            <h3 className="text-lg font-bold text-white mb-2">
              Confirm {confirmAction.action}
            </h3>
            <p className="text-sm text-[var(--color-muted)] mb-4">
              Are you sure you want to {confirmAction.action} container{" "}
              <strong className="text-white">
                {containers.find((c) => c.id === confirmAction.id)?.name}
              </strong>
              ?
            </p>
            <div className="flex gap-3 justify-end">
              <button
                onClick={() => setConfirmAction(null)}
                className="px-4 py-2 rounded-lg bg-[var(--color-background)] border border-[var(--color-border)] text-sm text-[var(--color-muted)] hover:text-white transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={() =>
                  executeAction(confirmAction.id, confirmAction.action)
                }
                className="px-4 py-2 rounded-lg bg-red-500/20 border border-red-500/30 text-sm text-red-400 hover:bg-red-500/30 transition-colors"
              >
                {confirmAction.action === "remove" ? "Remove" : "Stop"}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Logs Modal */}
      {showLogs && selectedContainer && (
        <div
          className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
          onClick={() => setShowLogs(false)}
        >
          <div
            className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl w-full max-w-3xl mx-4 max-h-[80vh] flex flex-col"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="flex items-center justify-between px-5 py-4 border-b border-[var(--color-border)]">
              <h3 className="text-sm font-bold text-white flex items-center gap-2">
                <ScrollText size={16} className="text-blue-400" />
                Logs — {selectedContainer.name}
              </h3>
              <button
                onClick={() => setShowLogs(false)}
                className="text-[var(--color-muted)] hover:text-white"
              >
                <X size={18} />
              </button>
            </div>
            <div className="flex-1 overflow-auto p-4 font-mono text-xs text-[var(--color-muted)] bg-[var(--color-background)] whitespace-pre-wrap">
              {logsContent}
            </div>
          </div>
        </div>
      )}

      {/* Inspect Modal */}
      {showInspect && selectedContainer && (
        <div
          className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
          onClick={() => setShowInspect(false)}
        >
          <div
            className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl w-full max-w-lg mx-4"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="flex items-center justify-between px-5 py-4 border-b border-[var(--color-border)]">
              <h3 className="text-sm font-bold text-white flex items-center gap-2">
                <Info size={16} className="text-purple-400" />
                Inspect — {selectedContainer.name}
              </h3>
              <button
                onClick={() => setShowInspect(false)}
                className="text-[var(--color-muted)] hover:text-white"
              >
                <X size={18} />
              </button>
            </div>
            <div className="p-5 space-y-3 font-mono text-xs">
              {[
                ["Name", selectedContainer.name],
                ["ID", selectedContainer.id],
                ["Image", selectedContainer.image],
                ["State", selectedContainer.state],
                ["Status", selectedContainer.status],
                ["Ports", selectedContainer.ports || "N/A"],
                ["Created", selectedContainer.created],
                [
                  "CPU",
                  selectedContainer.cpu !== undefined
                    ? `${selectedContainer.cpu.toFixed(1)}%`
                    : "N/A",
                ],
                [
                  "Memory",
                  selectedContainer.memory !== undefined
                    ? `${selectedContainer.memory.toFixed(1)}%`
                    : "N/A",
                ],
                ["Net In", selectedContainer.netIn || "N/A"],
                ["Net Out", selectedContainer.netOut || "N/A"],
              ].map(([k, v]) => (
                <div key={k} className="flex">
                  <span className="w-20 text-[var(--color-muted)] flex-shrink-0">
                    {k}:
                  </span>
                  <span className="text-white">{v}</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
