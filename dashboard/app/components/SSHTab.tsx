"use client";

import {
  AlertCircle,
  Check,
  Clock,
  Download,
  Key,
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
import { useI18n } from "../lib/i18n";
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
  const { t } = useI18n();
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
  const [showKeygen, setShowKeygen] = useState(false);
  const [keygenForm, setKeygenForm] = useState({
    name: "id_ed25519",
    type: "ed25519" as "ed25519" | "rsa",
    bits: 4096,
    comment: "",
    force: false,
  });
  const [keygenResult, setKeygenResult] = useState<{
    privateKey?: string;
    publicKey?: string;
  } | null>(null);
  const [keygenLoading, setKeygenLoading] = useState(false);

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

  async function generateKey() {
    setKeygenLoading(true);
    setKeygenResult(null);
    try {
      const result = await sshApi.keygen(keygenForm);
      if (result.success) {
        setKeygenResult({ privateKey: result.privateKey, publicKey: result.publicKey });
        showStatus(t("ssh.keygen.success"));
      } else if (result.exists) {
        // Key exists — show message and suggest enabling force
        showStatus(result.error || "Key đã tồn tại. Bật 'Ghi đè' để tạo mới.", true);
      } else {
        showStatus(result.error || t("ssh.keygen.error"), true);
      }
    } catch {
      showStatus(t("ssh.keygen.error"), true);
    } finally {
      setKeygenLoading(false);
    }
  }

  async function handleExport() {
    try {
      const hosts = await sshApi.export();
      const json = JSON.stringify(hosts, null, 2);
      const blob = new Blob([json], { type: "application/json" });
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = "ssh-hosts.json";
      a.click();
      URL.revokeObjectURL(url);
      showStatus(`${t("ssh.export")} OK — ${hosts.length} host(s)`);
    } catch {
      showStatus(t("ssh.export") + " thất bại", true);
    }
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
          <h2 className="text-2xl font-bold text-[var(--color-foreground)] flex items-center gap-3">
            <KeyRound size={24} className="text-[var(--success)]" />
            {t("ssh.title")}
          </h2>
          <p className="text-sm text-[var(--color-muted)] mt-1">
            {t("ssh.subtitle")}
          </p>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => testConnection()}
            disabled={testingAll}
            className="flex items-center gap-2 px-3 py-2 bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-muted)] hover:text-[var(--color-foreground)] hover:border-[var(--info)]/30 transition-all disabled:opacity-50"
          >
            <TestTube size={14} />
            {testingAll ? t("ssh.testing") : t("ssh.testAll")}
          </button>
          <button
            onClick={handleExport}
            className="flex items-center gap-2 px-3 py-2 bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-muted)] hover:text-[var(--color-foreground)] transition-all"
          >
            <Download size={14} />
            {t("ssh.export")}
          </button>
          <button className="flex items-center gap-2 px-3 py-2 bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-muted)] hover:text-[var(--color-foreground)] transition-all">
            <Upload size={14} />
            {t("ssh.import")}
          </button>
          <button
            onClick={() => { setShowKeygen(true); setKeygenResult(null); }}
            className="flex items-center gap-2 px-3 py-2 bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-muted)] hover:text-[var(--color-foreground)] hover:border-[var(--warning)]/30 transition-all"
          >
            <Key size={14} />
            {t("ssh.generateKey")}
          </button>
          <button
            onClick={openAddForm}
            className="flex items-center gap-2 px-3 py-2 bg-[var(--color-primary)] rounded-lg text-sm text-white hover:bg-[var(--color-primary-light)] transition-all"
          >
            <Plus size={14} />
            {t("ssh.addHost")}
          </button>
        </div>
      </div>

      {/* Stats */}
      <div className="flex gap-4">
        <div className="flex items-center gap-2 px-3 py-1.5 bg-[var(--color-card)] rounded-lg border border-[var(--color-border)]">
          <Server size={12} className="text-[var(--color-muted)]" />
          <span className="text-xs text-[var(--color-muted)]">
            {t("ssh.stats.total")}:
          </span>
          <span className="text-xs font-bold text-[var(--color-foreground)]">
            {hosts.length}
          </span>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-[var(--color-card)] rounded-lg border border-[var(--color-border)]">
          <Check size={12} className="text-[var(--success)]" />
          <span className="text-xs text-[var(--color-muted)]">
            {t("ssh.stats.connected")}:
          </span>
          <span className="text-xs font-bold text-[var(--success)]">
            {connectedCount}
          </span>
        </div>
        <div className="flex items-center gap-2 px-3 py-1.5 bg-[var(--color-card)] rounded-lg border border-[var(--color-border)]">
          <AlertCircle size={12} className="text-[var(--danger)]" />
          <span className="text-xs text-[var(--color-muted)]">
            {t("ssh.stats.failed")}:
          </span>
          <span className="text-xs font-bold text-[var(--danger)]">
            {failedCount}
          </span>
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
          placeholder={t("ssh.filterPlaceholder")}
          value={filter}
          onChange={(e) => setFilter(e.target.value)}
          className="w-full pl-10 pr-4 py-2.5 bg-[var(--color-card)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-foreground)] placeholder:text-[var(--color-muted)] focus:outline-none focus:border-[var(--color-primary)]"
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
                      ? "bg-[var(--success)]"
                      : host.status === "failed"
                        ? "bg-[var(--danger)]"
                        : "bg-[var(--color-muted)]"
                  }`}
                />
                <h3 className="font-bold text-[var(--color-foreground)] text-sm">
                  {host.name}
                </h3>
              </div>
              <span
                className={`text-[10px] px-2 py-0.5 rounded-full ${
                  host.status === "connected"
                    ? "bg-[var(--success)]/10 text-[var(--success)]"
                    : host.status === "failed"
                      ? "bg-[var(--danger)]/10 text-[var(--danger)]"
                      : "bg-[var(--color-muted)]/10 text-[var(--color-muted)]"
                }`}
              >
                {host.status}
              </span>
            </div>

            <div className="space-y-1.5 mb-4 font-mono text-xs">
              <div className="flex items-center gap-2">
                <span className="text-[var(--color-muted)] w-16">
                  {t("ssh.labels.host")}:
                </span>
                <span className="text-[var(--color-foreground)]">
                  {host.hostname}
                </span>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-[var(--color-muted)] w-16">
                  {t("ssh.labels.port")}:
                </span>
                <span className="text-[var(--color-foreground)]">
                  {host.port || "22"}
                </span>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-[var(--color-muted)] w-16">
                  {t("ssh.labels.user")}:
                </span>
                <span className="text-[var(--color-foreground)]">
                  {host.user || "-"}
                </span>
              </div>
              {host.identityFile && (
                <div className="flex items-center gap-2">
                  <span className="text-[var(--color-muted)] w-16">
                    {t("ssh.labels.key")}:
                  </span>
                  <span className="text-[var(--color-foreground)] truncate">
                    {host.identityFile}
                  </span>
                </div>
              )}
              {host.proxyJump && (
                <div className="flex items-center gap-2">
                  <span className="text-[var(--color-muted)] w-16">
                    {t("ssh.labels.proxy")}:
                  </span>
                  <span className="text-[var(--color-foreground)]">
                    {host.proxyJump}
                  </span>
                </div>
              )}
              {host.latency && (
                <div className="flex items-center gap-2">
                  <Clock size={10} className="text-[var(--color-muted)]" />
                  <span className="text-[var(--success)]">{host.latency}</span>
                </div>
              )}
            </div>

            <div className="flex items-center gap-1 pt-3 border-t border-[var(--color-border)]">
              <button
                onClick={() => connectToHost(host)}
                title={t("ssh.actions.connect")}
                className="flex-1 flex items-center justify-center gap-1.5 py-1.5 rounded-lg bg-[var(--success)]/10 text-[var(--success)] hover:bg-[var(--success)]/20 text-xs font-medium transition-colors"
              >
                <Plug size={12} /> {t("ssh.actions.connect")}
              </button>
              <button
                onClick={() => testConnection(host)}
                disabled={testingAll}
                title={t("ssh.actions.test")}
                className="p-1.5 rounded-lg hover:bg-[var(--info)]/10 text-[var(--color-muted)] hover:text-[var(--info)] transition-colors"
              >
                <TestTube size={14} />
              </button>
              <button
                onClick={() => openEditForm(host)}
                title={t("ssh.actions.edit")}
                className="p-1.5 rounded-lg hover:bg-[var(--warning)]/10 text-[var(--color-muted)] hover:text-[var(--warning)] transition-colors"
              >
                <Pencil size={14} />
              </button>
              <button
                onClick={() => setConfirmDelete(host.name)}
                title={t("ssh.actions.delete")}
                className="p-1.5 rounded-lg hover:bg-[var(--danger)]/10 text-[var(--color-muted)] hover:text-[var(--danger)] transition-colors"
              >
                <Trash2 size={14} />
              </button>
            </div>
          </div>
        ))}

        {filtered.length === 0 && (
          <div className="col-span-full text-center py-12 text-[var(--color-muted)]">
            <KeyRound size={32} className="mx-auto mb-3 opacity-50" />
            <p className="text-sm">{t("ssh.noHosts")}</p>
            <button
              onClick={openAddForm}
              className="mt-2 text-xs text-[var(--color-primary)] hover:text-[var(--color-primary-light)]"
            >
              {t("ssh.addFirst")}
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
              <h3 className="text-sm font-bold text-[var(--color-foreground)]">
                {editingHost
                  ? `${t("ssh.form.editTitle")} — ${editingHost.name}`
                  : t("ssh.form.addTitle")}
              </h3>
              <button
                onClick={() => setShowForm(false)}
                className="text-[var(--color-muted)] hover:text-[var(--color-foreground)]"
              >
                <X size={18} />
              </button>
            </div>
            <div className="p-5 space-y-3">
              {[
                {
                  label: t("ssh.form.name"),
                  key: "name" as const,
                  placeholder: "my-server",
                },
                {
                  label: t("ssh.form.hostname"),
                  key: "hostname" as const,
                  placeholder: "192.168.1.100",
                },
                {
                  label: t("ssh.form.port"),
                  key: "port" as const,
                  placeholder: "22",
                },
                {
                  label: t("ssh.form.user"),
                  key: "user" as const,
                  placeholder: "root",
                },
                {
                  label: t("ssh.form.identityFile"),
                  key: "identityFile" as const,
                  placeholder: "~/.ssh/id_rsa",
                },
                {
                  label: t("ssh.form.proxyJump"),
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
                    className="w-full px-3 py-2 bg-[var(--color-background)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-foreground)] placeholder:text-[var(--color-muted)] focus:outline-none focus:border-[var(--color-primary)] disabled:opacity-50"
                  />
                </div>
              ))}
            </div>
            <div className="flex gap-3 justify-end px-5 py-4 border-t border-[var(--color-border)]">
              <button
                onClick={() => setShowForm(false)}
                className="px-4 py-2 rounded-lg bg-[var(--color-background)] border border-[var(--color-border)] text-sm text-[var(--color-muted)] hover:text-[var(--color-foreground)] transition-colors"
              >
                {t("ssh.form.cancel")}
              </button>
              <button
                onClick={saveHost}
                className="px-4 py-2 rounded-lg bg-[var(--color-primary)] text-sm text-white hover:bg-[var(--color-primary-light)] transition-colors"
              >
                {editingHost ? t("ssh.form.update") : t("ssh.form.add")}
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
            <h3 className="text-lg font-bold text-[var(--color-foreground)] mb-2">
              {t("ssh.confirmDelete")}
            </h3>
            <p className="text-sm text-[var(--color-muted)] mb-4">
              Remove{" "}
              <strong className="text-[var(--color-foreground)]">
                {confirmDelete}
              </strong>{" "}
              {t("ssh.confirmDeleteDesc")}
            </p>
            <div className="flex gap-3 justify-end">
              <button
                onClick={() => setConfirmDelete(null)}
                className="px-4 py-2 rounded-lg bg-[var(--color-background)] border border-[var(--color-border)] text-sm text-[var(--color-muted)] hover:text-[var(--color-foreground)] transition-colors"
              >
                {t("common.cancel")}
              </button>
              <button
                onClick={() => deleteHost(confirmDelete)}
                className="px-4 py-2 rounded-lg bg-[var(--danger)]/20 border border-[var(--danger)]/30 text-sm text-[var(--danger)] hover:bg-[var(--danger)]/30 transition-colors"
              >
                {t("common.delete")}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Keygen Modal */}
      {showKeygen && (
        <div
          className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
          onClick={() => setShowKeygen(false)}
        >
          <div
            className="bg-[var(--color-card)] border border-[var(--color-border)] rounded-xl w-full max-w-md mx-4"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="flex items-center justify-between px-5 py-4 border-b border-[var(--color-border)]">
              <h3 className="text-sm font-bold text-[var(--color-foreground)] flex items-center gap-2">
                <Key size={16} className="text-[var(--warning)]" />
                {t("ssh.keygen.title")}
              </h3>
              <button
                onClick={() => setShowKeygen(false)}
                className="text-[var(--color-muted)] hover:text-[var(--color-foreground)]"
              >
                <X size={18} />
              </button>
            </div>
            <div className="p-5 space-y-3">
              <div>
                <label className="block text-xs text-[var(--color-muted)] mb-1">
                  {t("ssh.keygen.name")}
                </label>
                <input
                  type="text"
                  value={keygenForm.name}
                  onChange={(e) => setKeygenForm((f) => ({ ...f, name: e.target.value }))}
                  className="w-full px-3 py-2 bg-[var(--color-background)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-foreground)] focus:outline-none focus:border-[var(--color-primary)]"
                />
              </div>
              <div>
                <label className="block text-xs text-[var(--color-muted)] mb-1">
                  {t("ssh.keygen.type")}
                </label>
                <select
                  value={keygenForm.type}
                  onChange={(e) => setKeygenForm((f) => ({ ...f, type: e.target.value as "ed25519" | "rsa" }))}
                  className="w-full px-3 py-2 bg-[var(--color-background)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-foreground)] focus:outline-none focus:border-[var(--color-primary)]"
                >
                  <option value="ed25519">ed25519</option>
                  <option value="rsa">rsa</option>
                </select>
              </div>
              {keygenForm.type === "rsa" && (
                <div>
                  <label className="block text-xs text-[var(--color-muted)] mb-1">
                    {t("ssh.keygen.bits")}
                  </label>
                  <input
                    type="number"
                    value={keygenForm.bits}
                    onChange={(e) => setKeygenForm((f) => ({ ...f, bits: Number(e.target.value) }))}
                    className="w-full px-3 py-2 bg-[var(--color-background)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-foreground)] focus:outline-none focus:border-[var(--color-primary)]"
                  />
                </div>
              )}
              <div>
                <label className="block text-xs text-[var(--color-muted)] mb-1">
                  {t("ssh.keygen.comment")}
                </label>
                <input
                  type="text"
                  placeholder="you@example.com"
                  value={keygenForm.comment}
                  onChange={(e) => setKeygenForm((f) => ({ ...f, comment: e.target.value }))}
                  className="w-full px-3 py-2 bg-[var(--color-background)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-foreground)] placeholder:text-[var(--color-muted)] focus:outline-none focus:border-[var(--color-primary)]"
                />
              </div>
              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  id="keygen-force"
                  checked={keygenForm.force}
                  onChange={(e) => setKeygenForm((f) => ({ ...f, force: e.target.checked }))}
                  className="rounded border-[var(--color-border)]"
                />
                <label htmlFor="keygen-force" className="text-xs text-[var(--color-muted)]">
                  Ghi đè nếu key đã tồn tại (--force)
                </label>
              </div>
              {keygenResult && (
                <div className="bg-[var(--color-background)] rounded-lg p-3 space-y-1 font-mono text-xs">
                  <p className="text-[var(--color-muted)]">{t("ssh.keygen.privateKey")}:</p>
                  <p className="text-[var(--success)] break-all">{keygenResult.privateKey}</p>
                  <p className="text-[var(--color-muted)] mt-1">{t("ssh.keygen.publicKey")}:</p>
                  <p className="text-[var(--info)] break-all">{keygenResult.publicKey}</p>
                </div>
              )}
            </div>
            <div className="flex gap-3 justify-end px-5 py-4 border-t border-[var(--color-border)]">
              <button
                onClick={() => setShowKeygen(false)}
                className="px-4 py-2 rounded-lg bg-[var(--color-background)] border border-[var(--color-border)] text-sm text-[var(--color-muted)] hover:text-[var(--color-foreground)] transition-colors"
              >
                {t("ssh.form.cancel")}
              </button>
              <button
                onClick={generateKey}
                disabled={keygenLoading}
                className="px-4 py-2 rounded-lg bg-[var(--warning)]/20 border border-[var(--warning)]/30 text-sm text-[var(--warning)] hover:bg-[var(--warning)]/30 transition-colors disabled:opacity-50"
              >
                {keygenLoading ? t("ssh.keygen.generating") : t("ssh.keygen.generate")}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
