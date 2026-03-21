import { exec } from "child_process";
import { promisify } from "util";
import { NextResponse } from "next/server";

const execAsync = promisify(exec);
const CLI_PATH = process.env.IDOPS_CLI_PATH || "idops";

export async function GET(request: Request) {
  try {
    const { searchParams } = new URL(request.url);
    const protocol = searchParams.get("protocol") || "all";
    const portRange = searchParams.get("portRange") || "";

    let command = `${CLI_PATH} ports --json`;
    if (protocol && protocol !== "all") {
      command += ` --protocol ${protocol}`;
    }
    if (portRange) {
      command += ` --port ${portRange}`;
    }

    const { stdout } = await execAsync(command);
    const ports = JSON.parse(stdout);
    return NextResponse.json({ ports });
  } catch (error) {
    console.error("Ports scan error:", error);
    return NextResponse.json(
      { error: "Failed to scan ports", ports: [] },
      { status: 500 }
    );
  }
}
