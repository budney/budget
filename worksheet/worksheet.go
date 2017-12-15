package worksheet

import (
	"fmt"
	"github.com/araddon/dateparse"
	"google.golang.org/api/sheets/v4"
	"log"
	"time"
)

const HeaderRange = "B1:H1"
const DataRange = "B2:H"

