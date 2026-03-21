"use client";

import { useState } from "react";
import DockerTab from "./components/DockerTab";
import EnvTab from "./components/EnvTab";
import NginxTab from "./components/NginxTab";
import OverviewTab from "./components/OverviewTab";
import PortsTab from "./components/PortsTab";
import Sidebar, { type TabId } from "./components/Sidebar";
import SSHTab from "./components/SSHTab";
import {
  mockContainers,
  mockEnvVars,
  mockPorts,
  mockSSHHosts,
} from "./data/mockData";
import type { DockerContainer, EnvVariable, PortEntry, SSHHost } from "./types";

export default function Home() {
  const [activeTab, setActiveTab] = useState<TabId>("overview");
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

  // Initial mock data - components will fetch real data via APIs
  const [containers] = useState<DockerContainer[]>(mockContainers);
  const [ports] = useState<PortEntry[]>(mockPorts);
  const [sshHosts] = useState<SSHHost[]>(mockSSHHosts);
  const [envVars] = useState<EnvVariable[]>(mockEnvVars);

  return (
    <div className="min-h-screen bg-[var(--color-background)]">
      <Sidebar
        activeTab={activeTab}
        onTabChange={setActiveTab}
        collapsed={sidebarCollapsed}
        onToggle={() => setSidebarCollapsed(!sidebarCollapsed)}
      />

      <main
        className={`transition-all duration-300 min-h-screen ${
          sidebarCollapsed ? "ml-[68px]" : "ml-[240px]"
        }`}
      >
        <div className="p-6 max-w-[1400px]">
          {activeTab === "overview" && (
            <OverviewTab
              containers={containers}
              ports={ports}
              sshHosts={sshHosts}
              envVarCount={envVars.length}
              onNavigate={setActiveTab}
            />
          )}
          {activeTab === "docker" && <DockerTab containers={containers} />}
          {activeTab === "ports" && <PortsTab ports={ports} />}
          {activeTab === "ssh" && <SSHTab hosts={sshHosts} />}
          {activeTab === "env" && <EnvTab envVars={envVars} />}
          {activeTab === "nginx" && <NginxTab configs={[]} />}
        </div>
      </main>
    </div>
  );
}
