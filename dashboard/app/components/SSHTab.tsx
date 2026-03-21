"use client";

import {
  AlertCircle,
  Check,
  Clock,
  Download,
  KeyRound,
  Pencil,
  Plug,
  Plus,
  Server,
  TestTube,
  Trash2,
  Upload,
  X,
} from "lucide-react";
import { useEffect, useState } from "react";
import { sshApi } from "../lib/api";
import type { SSHHost } from "../types";

const emptyHost: SSHHost = {
  name: "",
  hostname: "",
  port: "22",
  user: "",
  identityFile: "",
  proxyJump: "",
  status: "unknown",
};

export default function SSHTab({ hosts: initialHosts }: { hosts: SSHHost[] }) {
  const [hosts, setHosts] = useState<SSHHost[]>(initialHosts);
  const [showForm, setShowForm] = useState(false);
  const [editingHost, setEditingHost] = useState<SSHHost | null>(null);
  const [formData, setFormData] = useState<SSHHost>({ ...emptyHost });
  const [confirmDelete, setConfirmDelete] = useState<string | null>(null);
  const [statusMsg, setStatusMsg] = useState<{
    text: string;
    isError: boolean;
  } | null>(null);
  const [testingAll, setTestingAll] = useState(false);
  const [filter, setFilter] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    loadHosts();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function loadHosts() {
    setIsLoading(true);
    try {
      const data = await sshApi.list();
      setHosts(data);
    } catch {
      showStatus("Failed to load SSH hosts", true);
    } finally {
      setIsLoading(false);
    }
  }

  async function testConnection(host?: SSHHost) {
    setTestingAll(true);
    try {
      const result = await sshApi.test(host?.name);
      if (result.results && result.results.length > 0) {
        const updatedHosts = hosts.map((h) => {
          const testResult = result.results.find((r) => r.Host.name === h.name);
          if (testResult) {
            const updatedHost: SSHHost = {
              ...h,
              status: (testResult.Success ? "connected" : "failed") as
                | "connected"
                | "failed",
            };
            if (testResult.Latency) updatedHost.latency = testResult.Latency;
            if (testResult.Error) updatedHost.error = testResult.Error;
            return updatedHost;
          }
          return h;
        });
        setHosts(updatedHosts);
        showStatus(host ? `Tested ${host.name}` : "All hosts tested");
      }
    } catch {
      showStatus("Connection test failed", true);
    } finally {
      setTestingAll(false);
    }
  }

  function showStatus(text: string, isError = false) {
    setStatusMsg({ text, isError });
    setTimeout(() => setStatusMsg(null), 3000);
  }

  function openAddForm() {
    setFormData({ ...emptyHost });
    setEditingHost(null);
    setShowForm(true);
  }

  function openEditForm(host: SSHHost) {
    setFormData({ ...host });
    setEditingHost(host);
    setShowForm(true);
  }

  function saveHost() {
    if (!formData.name || !formData.hostname) {
      showStatus("Name and Hostname are required", true);
      return;
    }

    if (editingHost) {
      setHosts((prev) =>
        prev.map((h) => (h.name === editingHost.name ? { ...formData } : h)),
      );
      showStatus(`Updated host ${formData.name}`);
    } else {
      if (hosts.some((h) => h.name === formData.name)) {
        showStatus(`Host "${formData.name}" already exists`, true);
        return;
      }
      setHosts((prev) => [...prev, { ...formData }]);
      showStatus(`Added host ${formData.name}`);
    }
    setShowForm(false);
  }

  function deleteHost(name: string) {
    setHosts((prev) => prev.filter((h) => h.name !== name));
    showStatus(`Deleted host ${name}`);
    setConfirmDelete(null);
  }

  function connectToHost(host: SSHHost) {
    showStatus(
      `Connecting to ${host.user || "root"}@${host.hostname}:${host.port || "22"} ...`,
    );
  }

  const filtered = hosts.filter(
    (h) =>
      h.name.toLowerCase().includes(filter.toLowerCase()) ||
      h.hostname.toLowerCase().includes(filter.toLowerCase()) ||
      h.user.toLowerCase().includes(filter.toLowerCase()),
  );

  const connectedCount = hosts.filter((h) => h.status === "connected").length;
  const failedCount = hosts.filter((h) => h.status === "failed").length;

  return (
    <div className="space-y-4 animate-fade-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-white flex items-center gap-3">
            <KeyRound size={24} className="text-emerald-400" />
            SSH Manager
          </h2>
          <p className="text-sm text-[var(--color-muted)] mt-1">
            Manage SSH hosts, connect, test connectivity, import/export config
          </p>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => testConnection()}
            disabled={testingAll}
            className="flex items-center gap-2 px-3 py-2 bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-muted)] hover:text-white hover:border-blue-500/30 transition-all disabled:opacity-50"
          >
            <TestTube size={14} />
            {testingAll ? "Testing..." : "Test All"}
          </button>
          <button className="flex items-center gap-2 px-3 py-2 bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-muted)] hover:text-white transition-all">
            <Download size={14} />
            Export
          </button>
          <button className="flex items-center gap-2 px-3 py-2 bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-muted)] hover:text-white transition-all">
            <Upload size={14} />
            Import
          </button>
          <button
            onClick={openAddForm}
            className="flex items-center gap-2 px-3 py-2 bg-[var(--color-primary)] rounded-lg text-sm text-white hover:bg-[var(--color-primary-light)] transition-all"
          >
            <Plus size={14} />
            Add Host
          </button>
        </div>
      </div>

      {/* Stats */}
      <div className="flex gap-4">
        <div className="flex items-center gap-2 px-3 py-1.5 bg-[var(--color-card)] rounded-lg border border-[var(--color-border)]">
          <Server size={12} className="text-[var(--color-muted)]" />
          <span className="text-xs text-[var(--color-muted)]">Total:</span>
          <span className="text-xs font-bold text-white">{hosts.length}</span>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-[var(--color-card)] rounded-lg border border-[var(--color-border)]">
          <Check size={12} className="text-emerald-400" />
          <span className="text-xs text-[var(--color-muted)]">Connected:</span>
          <span className="text-xs font-bold text-emerald-400">
            {connectedCount}
          </span>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-[var(--color-card)] rounded-lg border border-[var(--color-border)]">
          <AlertCircle size={12} className="text-red-400" />
          <span className="text-xs text-[var(--color-muted)]">Failed:</span>
          <span className="text-xs font-bold text-red-400">{failedCount}</span>
        </div>
      </div>

      {/* Status */}
      {statusMsg && (
        <div
          className={`px-4 py-2 rounded-lg text-sm ${statusMsg.isError ? "bg-red-500/10 text-red-400 border border-red-500/20" : "bg-emerald-500/10 text-emerald-400 border border-emerald-500/20"}`}
        >
          {statusMsg.text}
        </div>
      )}

      {/* Filter */}
      <div className="relative">
        <KeyRound
          size={16}
          className="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--color-muted)]"
        />
        <input
          type="text"
          placeholder="Filter hosts..."
          value={filter}
          onChange={(e) => setFilter(e.target.value)}
          className="w-full pl-10 pr-4 py-2.5 bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg text-sm text-white placeholder:text-[var(--color-muted)] focus:outline-none focus:border-[var(--color-primary)]"
        />
      </div>

      {/* Host Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
        {filtered.map((host) => (
          <div
            key={host.name}
            className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl p-4 hover:border-[var(--color-primary)]/30 transition-all group"
          >
            <div className="flex items-start justify-between mb-3">
              <div className="flex items-center gap-2">
                <div
                  className={`w-2.5 h-2.5 rounded-full ${
                    host.status === "connected"
                      ? "bg-emerald-400"
                      : host.status === "failed"
                        ? "bg-red-400"
                        : "bg-[var(--color-muted)]"
                  }`}
                />
                <h3 className="font-bold text-white text-sm">{host.name}</h3>
              </div>
              <span
                className={`text-[10px] px-2 py-0.5 rounded-full ${
                  host.status === "connected"
                    ? "bg-emerald-500/10 text-emerald-400"
                    : host.status === "failed"
                      ? "bg-red-500/10 text-red-400"
                      : "bg-[var(--color-muted)]/10 text-[var(--color-muted)]"
                }`}
              >
                {host.status}
              </span>
            </div>

            <div className="space-y-1.5 mb-4 font-mono text-xs">
              <div className="flex items-center gap-2">
                <span className="text-[var(--color-muted)] w-16">Host:</span>
                <span className="text-white">{host.hostname}</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-[var(--color-muted)] w-16">Port:</span>
                <span className="text-white">{host.port || "22"}</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-[var(--color-muted)] w-16">User:</span>
                <span className="text-white">{host.user || "-"}</span>
              </div>
              {host.identityFile && (
                <div className="flex items-center gap-2">
                  <span className="text-[var(--color-muted)] w-16">Key:</span>
                  <span className="text-white truncate">
                    {host.identityFile}
                  </span>
                </div>
              )}
              {host.proxyJump && (
                <div className="flex items-center gap-2">
                  <span className="text-[var(--color-muted)] w-16">Proxy:</span>
                  <span className="text-white">{host.proxyJump}</span>
                </div>
              )}
              {host.latency && (
                <div className="flex items-center gap-2">
                  <Clock size={10} className="text-[var(--color-muted)]" />
                  <span className="text-emerald-400">{host.latency}</span>
                </div>
              )}
            </div>

            <div className="flex items-center gap-1 pt-3 border-t border-[var(--color-border)]">
              <button
                onClick={() => connectToHost(host)}
                title="Connect"
                className="flex-1 flex items-center justify-center gap-1.5 py-1.5 rounded-lg bg-emerald-500/10 text-emerald-400 hover:bg-emerald-500/20 text-xs font-medium transition-colors"
              >
                <Plug size={12} /> Connect
              </button>
              <button
                onClick={() => testConnection(host)}
                disabled={testingAll}
                title="Test"
                className="p-1.5 rounded-lg hover:bg-blue-500/10 text-[var(--color-muted)] hover:text-blue-400 transition-colors"
              >
                <TestTube size={14} />
              </button>
              <button
                onClick={() => openEditForm(host)}
                title="Edit"
                className="p-1.5 rounded-lg hover:bg-amber-500/10 text-[var(--color-muted)] hover:text-amber-400 transition-colors"
              >
                <Pencil size={14} />
              </button>
              <button
                onClick={() => setConfirmDelete(host.name)}
                title="Delete"
                className="p-1.5 rounded-lg hover:bg-red-500/10 text-[var(--color-muted)] hover:text-red-400 transition-colors"
              >
                <Trash2 size={14} />
              </button>
            </div>
          </div>
        ))}

        {filtered.length === 0 && (
          <div className="col-span-full text-center py-12 text-[var(--color-muted)]">
            <KeyRound size={32} className="mx-auto mb-3 opacity-50" />
            <p className="text-sm">No SSH hosts found</p>
            <button
              onClick={openAddForm}
              className="mt-2 text-xs text-[var(--color-primary)] hover:text-[var(--color-primary-light)]"
            >
              Add your first host
            </button>
          </div>
        )}
      </div>

      {/* Add/Edit Form Modal */}
      {showForm && (
        <div
          className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
          onClick={() => setShowForm(false)}
        >
          <div
            className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl w-full max-w-md mx-4"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="flex items-center justify-between px-5 py-4 border-b border-[var(--color-border)]">
              <h3 className="text-sm font-bold text-white">
                {editingHost
                  ? `Edit Host — ${editingHost.name}`
                  : "Add New Host"}
              </h3>
              <button
                onClick={() => setShowForm(false)}
                className="text-[var(--color-muted)] hover:text-white"
              >
                <X size={18} />
              </button>
            </div>
            <div className="p-5 space-y-3">
              {[
                {
                  label: "Host Name *",
                  key: "name" as const,
                  placeholder: "my-server",
                },
                {
                  label: "Hostname *",
                  key: "hostname" as const,
                  placeholder: "192.168.1.100",
                },
                { label: "Port", key: "port" as const, placeholder: "22" },
                { label: "User", key: "user" as const, placeholder: "root" },
                {
                  label: "Identity File",
                  key: "identityFile" as const,
                  placeholder: "~/.ssh/id_rsa",
                },
                {
                  label: "ProxyJump",
                  key: "proxyJump" as const,
                  placeholder: "bastion",
                },
              ].map((field) => (
                <div key={field.key}>
                  <label className="block text-xs text-[var(--color-muted)] mb-1">
                    {field.label}
                  </label>
                  <input
                    type="text"
                    placeholder={field.placeholder}
                    value={formData[field.key]}
                    onChange={(e) =>
                      setFormData((prev) => ({
                        ...prev,
                        [field.key]: e.target.value,
                      }))
                    }
                    disabled={editingHost !== null && field.key === "name"}
                    className="w-full px-3 py-2 bg-[var(--color-background)] border border-[var(--color-border)] rounded-lg text-sm text-white placeholder:text-[var(--color-muted)] focus:outline-none focus:border-[var(--color-primary)] disabled:opacity-50"
                  />
                </div>
              ))}
            </div>
            <div className="flex gap-3 justify-end px-5 py-4 border-t border-[var(--color-border)]">
              <button
                onClick={() => setShowForm(false)}
                className="px-4 py-2 rounded-lg bg-[var(--color-background)] border border-[var(--color-border)] text-sm text-[var(--color-muted)] hover:text-white transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={saveHost}
                className="px-4 py-2 rounded-lg bg-[var(--color-primary)] text-sm text-white hover:bg-[var(--color-primary-light)] transition-colors"
              >
                {editingHost ? "Update" : "Add Host"}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Confirm Delete */}
      {confirmDelete && (
        <div
          className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
          onClick={() => setConfirmDelete(null)}
        >
          <div
            className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl p-6 max-w-sm w-full mx-4"
            onClick={(e) => e.stopPropagation()}
          >
            <h3 className="text-lg font-bold text-white mb-2">Delete Host</h3>
            <p className="text-sm text-[var(--color-muted)] mb-4">
              Remove <strong className="text-white">{confirmDelete}</strong>{" "}
              from SSH config? A backup will be created automatically.
            </p>
            <div className="flex gap-3 justify-end">
              <button
                onClick={() => setConfirmDelete(null)}
                className="px-4 py-2 rounded-lg bg-[var(--color-background)] border border-[var(--color-border)] text-sm text-[var(--color-muted)] hover:text-white transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={() => deleteHost(confirmDelete)}
                className="px-4 py-2 rounded-lg bg-red-500/20 border border-red-500/30 text-sm text-red-400 hover:bg-red-500/30 transition-colors"
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
