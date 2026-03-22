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
  } catch {
    return NextResponse.json({ hosts: [], unavailable: true });
  }
}
