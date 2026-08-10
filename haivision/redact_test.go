package haivision

import (
	"net/http"
	"strings"
	"testing"
)

// Test in-package: redactJSON e redactHeader sono non esportate. È l'unica ragione per cui
// questo file non sta in test/ come gli altri.

func TestRedactJSONMascheraLeCredenziali(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		absent  []string // sottostringhe che NON devono comparire
		present []string // sottostringhe che devono comparire
	}{
		{
			name:    "body di login",
			in:      `{"username":"haiadmin","password":"sup3r-s3cret"}`,
			absent:  []string{"sup3r-s3cret"},
			present: []string{`"username":"haiadmin"`, `"password":"[REDACTED]"`},
		},
		{
			name: "risposta di sessione",
			in: `{"response":{"type":"Session","sessionID":"abc123XYZ",` +
				`"numLoginFailures":0}}`,
			absent:  []string{"abc123XYZ"},
			present: []string{`"sessionID":"[REDACTED]"`, `"numLoginFailures":0`},
		},
		{
			name:    "passphrase SRT",
			in:      `{"name":"dst","srtPassPhrase":"my-passphrase","port":2000}`,
			absent:  []string{"my-passphrase"},
			present: []string{`"srtPassPhrase":"[REDACTED]"`, `"port":2000`},
		},
		{
			name:    "chiave con case diverso e spazi",
			in:      `{"Password" : "hunter2"}`,
			absent:  []string{"hunter2"},
			present: []string{"[REDACTED]"},
		},
		{
			name:    "valore con virgolette escapate",
			in:      `{"password":"a\"b","keep":"me"}`,
			absent:  []string{`a\"b`},
			present: []string{`"password":"[REDACTED]"`, `"keep":"me"`},
		},
		{
			name:    "niente da mascherare resta identico",
			in:      `{"routeID":"route-1","state":"connected"}`,
			present: []string{`"routeID":"route-1"`, `"state":"connected"`},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := redactJSON(tc.in)
			for _, s := range tc.absent {
				if strings.Contains(got, s) {
					t.Errorf("il valore sensibile %q è ancora presente: %s", s, got)
				}
			}
			for _, s := range tc.present {
				if !strings.Contains(got, s) {
					t.Errorf("atteso %q nel risultato, ottenuto: %s", s, got)
				}
			}
		})
	}
}

func TestRedactJSONBodyVuoto(t *testing.T) {
	if got := redactJSON(""); got != "" {
		t.Errorf("redactJSON(\"\") = %q, atteso \"\"", got)
	}
}

func TestRedactHeaderMascheraSenzaMutareLInput(t *testing.T) {
	in := http.Header{}
	in.Set("Authorization", "Basic dXNlcjpwYXNz")
	in.Set("Cookie", "sessionID=abc123XYZ")
	in.Set("X-Tenant", "hiway")
	in.Add("Accept", "application/json")

	out := redactHeader(in)

	if got := out.Get("Authorization"); got != redactedPlaceholder {
		t.Errorf("Authorization = %q, atteso %s", got, redactedPlaceholder)
	}
	if got := out.Get("Cookie"); got != redactedPlaceholder {
		t.Errorf("Cookie = %q, atteso %s (il cookie sessionID è una credenziale)", got, redactedPlaceholder)
	}
	if got := out.Get("X-Tenant"); got != "hiway" {
		t.Errorf("X-Tenant = %q, atteso hiway (non è sensibile)", got)
	}
	if got := out.Get("Accept"); got != "application/json" {
		t.Errorf("Accept = %q", got)
	}
	// l'input non deve essere toccato: gli header di una richiesta in volo non vanno mutati
	if got := in.Get("Authorization"); got != "Basic dXNlcjpwYXNz" {
		t.Errorf("redactHeader ha mutato l'input: Authorization = %q", got)
	}
	if got := in.Get("Cookie"); got != "sessionID=abc123XYZ" {
		t.Errorf("redactHeader ha mutato l'input: Cookie = %q", got)
	}
}

func TestBodyExcerptTroncaERedige(t *testing.T) {
	// il body d'errore di un gateway può essere una pagina HTML intera
	long := strings.Repeat("à", bodyExcerptLimit+50) // multi-byte: verifica il taglio per rune
	got := bodyExcerpt([]byte(long))
	if !strings.HasSuffix(got, "… (troncato)") {
		t.Errorf("atteso suffisso di troncamento, ottenuto: %q", got[max(0, len(got)-30):])
	}
	if n := len([]rune(strings.TrimSuffix(got, "… (troncato)"))); n != bodyExcerptLimit {
		t.Errorf("troncato a %d rune, atteso %d", n, bodyExcerptLimit)
	}
	if !strings.Contains(bodyExcerpt([]byte(`{"password":"x"}`)), redactedPlaceholder) {
		t.Error("bodyExcerpt non ha redatto la password")
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
