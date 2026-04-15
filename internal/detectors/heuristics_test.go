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
	src := `int main(){int x1=0; int y2=1; return x1+y2;}`
	outs := Run(ingestion.Prepare(src))
	if outs.DeadCode.Score != 0 || outs.ControlFlow.Score != 0 {
		t.Fatalf("expected neutral non-identifier detectors, got %+v", outs)
	}
}
