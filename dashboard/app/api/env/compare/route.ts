import { exec } from "child_process";
import { promisify } from "util";
import { NextResponse } from "next/server";

const execAsync = promisify(exec);
const CLI_PATH = process.env.IDOPS_CLI_PATH || "idops";

export async function GET(request: Request) {
  try {
    const { searchParams } = new URL(request.url);
    const source = searchParams.get("source") || ".env.example";
    const target = searchParams.get("target") || ".env";

    const { stdout } = await execAsync(
      `${CLI_PATH} env compare --source ${source} --target ${target}`
    );
    return NextResponse.json({ output: stdout });
  } catch (error) {
    console.error("Env compare error:", error);
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    return NextResponse.json(
      { error: errorMessage, output: "" },
      { status: 500 }
    );
  }
}
