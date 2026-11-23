import {
  CheckIcon,
  CircleIcon,
  ShieldCheckIcon,
  UploadIcon,
} from "lucide-react";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router";

import Container from "src/components/Container";
import TwoColumnLayoutHeader from "src/components/TwoColumnLayoutHeader";
import { Alert, AlertDescription, AlertTitle } from "src/components/ui/alert";
import { Button } from "src/components/ui/button";
import { Label } from "src/components/ui/label";
import { Progress } from "src/components/ui/progress";
import { toast } from "sonner";
import useSetupStore from "src/state/SetupStore";
import { request } from "src/utils/request";

type CertStatus = {
  stage: "idle" | "uploading" | "validating" | "storing" | "ready" | "failed";
  progress: number;
  message: string;
  error?: string;
};

const steps = [
  { label: "Upload certificate", at: 10 },
  { label: "Validate certificate and key", at: 35 },
  { label: "Store encrypted on this Hub", at: 70 },
  { label: "Ready to continue", at: 100 },
];

export function SetupGreenlightCert() {
  const navigate = useNavigate();
  const unlockPassword = useSetupStore((state) => state.unlockPassword);
  const [file, setFile] = useState<File | null>(null);
  const [status, setStatus] = useState<CertStatus>({
    stage: "idle",
    progress: 0,
    message: "Waiting for certificate",
  });
  const [waiting, setWaiting] = useState(false);
  const [uploading, setUploading] = useState(false);

  useEffect(() => {
    (async () => {
      try {
        const current = await request<CertStatus>(
          "/api/setup/greenlight-cert/status"
        );
        if (current) {
          setStatus(current);
        }
      } catch {
        // status is optional until the first upload
      }
    })();
  }, []);

  useEffect(() => {
    if (!waiting) {
      return;
    }
    let active = true;
    const poll = async () => {
      try {
        const current = await request<CertStatus>(
          "/api/setup/greenlight-cert/status"
        );
        if (!active || !current) {
          return;
        }
        setStatus(current);
        if (current.stage === "ready" || current.stage === "failed") {
          setWaiting(false);
        }
      } catch {
        // keep polling; the upload may still be starting
      }
    };
    const interval = setInterval(poll, 400);
    poll();
    return () => {
      active = false;
      clearInterval(interval);
    };
  }, [waiting]);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!file) {
      toast.error("Choose the certificate zip first");
      return;
    }
    if (!unlockPassword) {
      toast.error("Unlock password is missing");
      return;
    }

    const form = new FormData();
    form.set("unlockPassword", unlockPassword);
    form.set("file", file);

    try {
      setUploading(true);
      await request("/api/setup/greenlight-cert", {
        method: "POST",
        body: form,
      });
      setStatus({
        stage: "uploading",
        progress: 10,
        message: "Uploading certificate",
      });
      setWaiting(true);
    } catch (error) {
      toast.error("Upload failed", {
        description: error instanceof Error ? error.message : undefined,
      });
    } finally {
      setUploading(false);
    }
  }

  const ready = status.stage === "ready" && status.progress === 100;
  const failed = status.stage === "failed";

  return (
    <Container>
      <TwoColumnLayoutHeader
        title="Greenlight certificate"
        pageTitle="Greenlight certificate"
        description="Upload the zip from the Greenlight Developer Console. We use it only to register or recover this node."
      />

      <form onSubmit={onSubmit} className="mt-6 flex max-w-xl flex-col gap-6">
        <div className="grid gap-2">
          <Label htmlFor="greenlight-cert-zip">Certificate zip</Label>
          <input
            id="greenlight-cert-zip"
            type="file"
            accept=".zip,application/zip"
            disabled={uploading || waiting}
            onChange={(e) => setFile(e.target.files?.[0] ?? null)}
            className="text-sm file:mr-3 file:rounded-md file:border file:px-3 file:py-1 file:text-sm"
          />
          <p className="text-sm text-muted-foreground">
            Must contain client.crt and client-key.pem.
          </p>
        </div>

        <Button type="submit" disabled={!file || uploading || waiting}>
          <UploadIcon className="size-4" />
          {uploading ? "Uploading..." : "Upload & verify"}
        </Button>

        {(waiting || status.stage !== "idle") && (
          <div className="rounded-lg border p-4">
            <div className="mb-2 flex items-center justify-between text-sm">
              <span>{status.message}</span>
              <span>{status.progress}%</span>
            </div>
            <Progress value={status.progress} />
            <div className="mt-4 space-y-2">
              {steps.map((step) => {
                const done = status.progress >= step.at;
                return (
                  <div
                    key={step.label}
                    className="flex items-center gap-2 text-sm"
                  >
                    {done ? (
                      <CheckIcon className="size-4 text-primary" />
                    ) : (
                      <CircleIcon className="size-4 text-muted-foreground" />
                    )}
                    <span
                      className={
                        done ? "text-foreground" : "text-muted-foreground"
                      }
                    >
                      {step.label}
                    </span>
                  </div>
                );
              })}
            </div>
          </div>
        )}

        {failed && (
          <Alert variant="destructive">
            <AlertTitle>Certificate not accepted</AlertTitle>
            <AlertDescription>{status.error}</AlertDescription>
          </Alert>
        )}

        {ready && (
          <Alert>
            <ShieldCheckIcon className="size-4" />
            <AlertTitle>Certificate ready</AlertTitle>
            <AlertDescription>
              The developer certificate is stored encrypted on this Hub.
            </AlertDescription>
          </Alert>
        )}

        <Button
          type="button"
          disabled={!ready}
          onClick={() => navigate("/setup/security")}
        >
          Resume setup
        </Button>
      </form>
    </Container>
  );
}
