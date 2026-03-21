import { exec } from "child_process";
import { NextResponse } from "next/server";
import { promisify } from "util";

const execAsync = promisify(exec);
const CLI_PATH = process.env.IDOPS_CLI_PATH || "idops";

export async function GET(request: Request) {
  try {
    const { searchParams } = new URL(request.url);
    const file = searchParams.get("file") || ".env";

    const { stdout } = await execAsync(
      `${CLI_PATH} env validate --file ${file}`,
    );
    return NextResponse.json({ valid: true, output: stdout });
  } catch (error) {
    console.error("Env validate error:", error);
    const errorMessage =
      error instanceof Error ? error.message : "Validation failed";
    return NextResponse.json(
      { valid: false, error: errorMessage, output: "" },
      { status: 200 },
    );
  }
}
