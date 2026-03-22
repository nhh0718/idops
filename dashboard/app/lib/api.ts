import { DockerContainer, PortEntry, SSHHost } from "../types";

// API client for idops dashboard
const API_BASE = "/api";

// Docker API
export const dockerApi = {
  list: async (): Promise<DockerContainer[]> => {
    const res = await fetch(`${API_BASE}/docker`);
    const data = await res.json();
    return data.containers || [];
  },

  action: async (
    action: string,
    containerId: string,
  ): Promise<{ success: boolean; message?: string; error?: string }> => {
    const res = await fetch(`${API_BASE}/docker/action`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ action, containerId }),
    });
    return res.json();
  },

  logs: async (containerId: string): Promise<string> => {
    const res = await fetch(
      `${API_BASE}/docker/logs?containerId=${containerId}`,
    );
    const data = await res.json();
    return data.logs || "";
  },
};

// Ports API
export const portsApi = {
  scan: async (protocol?: string, portRange?: string): Promise<PortEntry[]> => {
    const params = new URLSearchParams();
    if (protocol && protocol !== "all") params.append("protocol", protocol);
    if (portRange) params.append("portRange", portRange);

    const res = await fetch(`${API_BASE}/ports?${params}`);
    const data = await res.json();
    return data.ports || [];
  },

  kill: async (
    port: number,
  ): Promise<{ success: boolean; message?: string; error?: string }> => {
    const res = await fetch(`${API_BASE}/ports/kill`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ port }),
    });
    return res.json();
  },
};

// SSH API
export const sshApi = {
  list: async (): Promise<SSHHost[]> => {
    const res = await fetch(`${API_BASE}/ssh`);
    const data = await res.json();
    return data.hosts || [];
  },

  test: async (
    host?: string,
  ): Promise<{
    results: {
      Host: SSHHost;
      Success: boolean;
      Latency: string;
      Error?: string;
    }[];
    unavailable?: boolean;
  }> => {
    const params = host ? `?host=${host}` : "";
    const res = await fetch(`${API_BASE}/ssh/test${params}`);
    return res.json();
  },

  keygen: async (opts: {
    name?: string;
    type?: "ed25519" | "rsa";
    bits?: number;
    comment?: string;
  }): Promise<{ success: boolean; privateKey?: string; publicKey?: string; error?: string }> => {
    const res = await fetch(`${API_BASE}/ssh/keygen`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(opts),
    });
    return res.json();
  },

  export: async (): Promise<SSHHost[]> => {
    const res = await fetch(`${API_BASE}/ssh/export`);
    const data = await res.json();
    return data.hosts || [];
  },
};

// Env API
export const envApi = {
  show: async (file = ".env"): Promise<Record<string, string>> => {
    const res = await fetch(`${API_BASE}/env?file=${file}`);
    const data = await res.json();
    return data.envVars || {};
  },

  compare: async (
    source = ".env.example",
    target = ".env",
  ): Promise<{ output: string }> => {
    const res = await fetch(
      `${API_BASE}/env/compare?source=${source}&target=${target}`,
    );
    return res.json();
  },

  validate: async (
    file = ".env",
  ): Promise<{ valid: boolean; output?: string; error?: string }> => {
    const res = await fetch(`${API_BASE}/env/validate?file=${file}`);
    return res.json();
  },
};

// Nginx API
export const nginxApi = {
  list: async (): Promise<string[]> => {
    const res = await fetch(`${API_BASE}/nginx`);
    const data = await res.json();
    return data.configs || [];
  },

  validate: async (): Promise<{
    valid: boolean;
    message?: string;
    error?: string;
  }> => {
    const res = await fetch(`${API_BASE}/nginx/validate`);
    return res.json();
  },
};
