import { exec } from "child_process";
import { promisify } from "util";
import { NextResponse } from "next/server";

const execAsync = promisify(exec);
const CLI_PATH = process.env.IDOPS_CLI_PATH || "idops";

export async function GET(request: Request) {
  try {
    const { searchParams } = new URL(request.url);
    const host = searchParams.get("host");

    let command = `${CLI_PATH} ssh test --json`;
    if (host) {
      command += ` ${host}`;
    }

    const { stdout } = await execAsync(command);
    const results = JSON.parse(stdout);
    return NextResponse.json({ results });
  } catch (error) {
    console.error("SSH test error:", error);
    return NextResponse.json({ results: [], unavailable: true });
  }
}
