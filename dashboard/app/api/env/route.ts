import { exec } from "child_process";
import { promisify } from "util";
import { NextResponse } from "next/server";

const execAsync = promisify(exec);
const CLI_PATH = process.env.IDOPS_CLI_PATH || "idops";

export async function GET(request: Request) {
  try {
    const { searchParams } = new URL(request.url);
    const file = searchParams.get("file") || ".env";

    const { stdout } = await execAsync(`${CLI_PATH} env show --file ${file} --json`);
    const envVars = JSON.parse(stdout);
    return NextResponse.json({ envVars });
  } catch (error) {
    console.error("Env show error:", error);
    return NextResponse.json(
      { error: "Failed to read env file", envVars: {} },
      { status: 500 }
    );
  }
}
