import { exec } from "child_process";
import { promisify } from "util";
import { NextResponse } from "next/server";

const execAsync = promisify(exec);
const CLI_PATH = process.env.IDOPS_CLI_PATH || "idops";

export async function GET() {
  try {
    const { stdout } = await execAsync(`${CLI_PATH} nginx validate`);
    return NextResponse.json({ valid: true, message: stdout });
  } catch (error) {
    console.error("Nginx validate error:", error);
    const errorMessage = error instanceof Error ? error.message : "Validation failed";
    return NextResponse.json(
      { valid: false, error: errorMessage },
      { status: 200 }
    );
  }
}
