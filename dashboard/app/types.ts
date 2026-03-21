export interface DockerContainer {
  id: string;
  name: string;
  image: string;
  state: "running" | "exited" | "paused" | "restarting" | "created";
  status: string;
  ports: string;
  cpu?: number;
  memory?: number;
  memUsage?: string;
  memLimit?: string;
  netIn?: string;
  netOut?: string;
  created: string;
}

export interface PortEntry {
  protocol: string;
  address: string;
  port: number;
  pid: number;
  process: string;
  user: string;
  status: string;
}

export interface SSHHost {
  name: string;
  hostname: string;
  port: string;
  user: string;
  identityFile: string;
  proxyJump: string;
  status: "connected" | "failed" | "unknown";
  latency?: string;
  error?: string;
}

export interface EnvVariable {
  key: string;
  value: string;
  isSensitive: boolean;
  comment?: string;
}

export interface EnvFile {
  path: string;
  variables: EnvVariable[];
}

export interface EnvDiff {
  missing: string[];
  extra: string[];
  changed: { key: string; source: string; target: string }[];
}

export interface EnvValidationIssue {
  line: number;
  key: string;
  type:
    | "empty"
    | "duplicate"
    | "trailing_space"
    | "invalid_format"
    | "unquoted_spaces";
  message: string;
}

export interface NginxConfig {
  template: string;
  serverName: string;
  listenPort: number;
  sslEnabled: boolean;
  sslCertPath: string;
  sslKeyPath: string;
  upstreamHost?: string;
  upstreamPort?: number;
  webSocket?: boolean;
  rootPath?: string;
  phpSocket?: string;
  backends?: { host: string; port: number; weight: number }[];
  method?: string;
  upstreamName?: string;
  enableGzip?: boolean;
  cacheMaxAge?: number;
}
