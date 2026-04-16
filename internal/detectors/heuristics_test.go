package detectors

import (
	"testing"

	"github.com/graphsentinel/graphsentinel/internal/ingestion"
)

func TestDetectIdentifierRenaming_Obfuscated(t *testing.T) {
	t.Parallel()
	d := HeuristicIdentifierRenamingDetector{}
	src := `
int main(){
  int a1=0;
  int bb=3;
  int zzzzzz=4;
  int x9qk2m=1;
  int r4t7p1=2;
  return a1+bb+zzzzzz+x9qk2m+r4t7p1;
}`
	out := d.Detect(ingestion.Prepare(src))
	if !out.Likely {
		t.Fatalf("expected likely obfuscated naming, got %+v", out)
	}
}

func TestDetectIdentifierRenaming_ClearNames(t *testing.T) {
	t.Parallel()
	d := HeuristicIdentifierRenamingDetector{}
	src := `
int compute_total(int order_count, int discount_rate){
  int total_value = order_count * 5;
  return total_value - discount_rate;
}`
	out := d.Detect(ingestion.Prepare(src))
	if out.Likely {
		t.Fatalf("expected normal naming, got %+v", out)
	}
}

func TestRun_OnlyIdentifierActiveForNow(t *testing.T) {
	t.Parallel()
	src := `int main(){int x1=0; int y2=1; if(false){return 1;} return x1+y2;}`
	outs := Run(ingestion.Prepare(src))
	if !outs.DeadCode.Likely || outs.DeadCode.Score <= 0 {
		t.Fatalf("expected dead-code signal to be active, got %+v", outs.DeadCode)
	}
	if outs.ControlFlow.Score != 0 || outs.ControlFlow.Likely {
		t.Fatalf("expected neutral control-flow detector, got %+v", outs.ControlFlow)
	}
}

func TestDetectDeadCode_HeuristicDetector(t *testing.T) {
	t.Parallel()
	d := HeuristicDeadCodeDetector{}

	obfuscated := `
int main() {
  int unused_tmp = 0;
  if(false) { return 1; }
  while(0) { unused_tmp++; }
  return 0;
}`
	out := d.Detect(ingestion.Prepare(obfuscated))
	if !out.Likely {
		t.Fatalf("expected dead-code detection, got %+v", out)
	}

	clear := `
int main() {
  int count = 2;
  if (count > 0) { return count; }
  return 0;
}`
	out2 := d.Detect(ingestion.Prepare(clear))
	if out2.Likely {
		t.Fatalf("expected low dead-code signal, got %+v", out2)
	}
}
