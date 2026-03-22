import { exec } from "child_process";
import { promisify } from "util";
import { NextResponse } from "next/server";

const execAsync = promisify(exec);
const CLI_PATH = process.env.IDOPS_CLI_PATH || "idops";

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const { name = "id_ed25519", type = "ed25519", bits = 4096, comment = "", force = false } = body;

    if (type !== "ed25519" && type !== "rsa") {
      return NextResponse.json(
        { error: "Loại key không hợp lệ, chỉ hỗ trợ ed25519 hoặc rsa" },
        { status: 400 }
      );
    }

    let command = `${CLI_PATH} ssh keygen --json --name ${name} --type ${type}`;
    if (force) {
      command += " --force";
    }
    if (type === "rsa") {
      command += ` --bits ${bits}`;
    }
    if (comment) {
      command += ` --comment "${comment}"`;
    }

    const { stdout } = await execAsync(command);
    const result = JSON.parse(stdout);
    return NextResponse.json({ success: true, ...result });
  } catch (error: unknown) {
    const stderr = (error as { stderr?: string })?.stderr || "";
    // Key already exists — return as info, not server error
    if (stderr.includes("đã tồn tại")) {
      return NextResponse.json({
        success: false,
        error: "Key đã tồn tại. Bật 'Ghi đè' để tạo mới.",
        exists: true,
      });
    }
    return NextResponse.json({
      error: "Tạo SSH key thất bại",
      success: false,
    });
  }
}
