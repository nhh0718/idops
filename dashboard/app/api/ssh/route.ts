import { exec } from "child_process";
import { promisify } from "util";
import { NextResponse } from "next/server";

const execAsync = promisify(exec);
const CLI_PATH = process.env.IDOPS_CLI_PATH || "idops";

export async function GET() {
  try {
    const { stdout } = await execAsync(`${CLI_PATH} ssh --json`);
    const hosts = JSON.parse(stdout);
    return NextResponse.json({ hosts });
  } catch (error) {
    console.error("SSH list error:", error);
    return NextResponse.json(
      { error: "Failed to list SSH hosts", hosts: [] },
      { status: 500 }
    );
  }
}
