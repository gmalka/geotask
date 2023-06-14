package geo

import (
	"math"
	"math/rand"

	geo "github.com/kellydunn/golang-geo"
)

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

//go:generate mockery --name PolygonChecker
type PolygonChecker interface {
	Contains(point Point) bool // проверить, находится ли точка внутри полигона
	Allowed() bool             // разрешено ли входить в полигон
	RandomPoint() Point        // сгенерировать случайную точку внутри полигона
}

type Polygon struct {
	polygon *geo.Polygon
	allowed bool
}

func NewPolygon(points []Point, allowed bool) *Polygon {
	// используем библиотеку golang-geo для создания полигона
	geoPoints := make([]*geo.Point, len(points))
	for i := 0; i < len(points); i++ {
		geoPoints[i] = geo.NewPoint(points[i].Lat, points[i].Lng)
	}

	return &Polygon{
		polygon: geo.NewPolygon(geoPoints),
		allowed: allowed,
	}
}

func (p *Polygon) Contains(point Point) bool {
	return p.polygon.Contains(geo.NewPoint(point.Lat, point.Lng))
}

func (p *Polygon) Allowed() bool {
	return p.allowed
}

func (p *Polygon) RandomPoint() Point {
	points := p.polygon.Points()
	for {
		var maxLat, maxLng float64
		minLat := math.MaxFloat64
		minLng := math.MaxFloat64
		for i := 0; i < len(points); i++ {
			if minLat > points[i].Lat() {
				minLat = points[i].Lat()
			}
			if points[i].Lat() > maxLat {
				maxLat = points[i].Lat()
			}

			if minLng > points[i].Lng() {
				minLng = points[i].Lng()
			}
			if points[i].Lng() > maxLng {
				maxLng = points[i].Lng()
			}
		}

		randLat := minLat + rand.Float64() * (maxLat - minLat)
		randLng := minLng + rand.Float64() * (maxLng - minLng)
		
		point := geo.NewPoint(randLat, randLng)
		if p.polygon.Contains(point) {
			return Point{randLat, randLng}
		}
	}
}

func CheckPointIsAllowed(point Point, allowedZone PolygonChecker, disabledZones []PolygonChecker) bool {
	// проверить, находится ли точка в разрешенной зоне

	if allowedZone.Contains(point) {
		for _, v := range disabledZones {
			if v.Contains(point) {
				return false
			}
		}
	} else {
		return false
	}

	return true
}

func GetRandomAllowedLocation(allowedZone PolygonChecker, disabledZones []PolygonChecker) Point {
	var point Point
	// получение случайной точки в разрешенной зоне
	for {
		point = allowedZone.RandomPoint()
		for i := 0; i < len(disabledZones) + 1; i++ {
			if i == len(disabledZones) {
				return point
			}
			if disabledZones[i].Contains(point) {
				break
			}
		}
	}
}

func NewDisAllowedZone1() *Polygon {
	points := []Point{{60.051063834232714, 30.28244720269174},
		{60.0509781359604, 30.341498716363613},
		{60.02036963316746, 30.363471372613613},
		{60.01650940538451, 30.31986938286752},}

	return NewPolygon(points, false)
}

func NewDisAllowedZone2() *Polygon {
	points := []Point{{59.902742187627325, 30.35368172093575},
		{59.90015959974209, 30.41290489598458},
		{59.842429456164574, 30.411531604968953},
		{59.836047143247896, 30.373766102039266},}

	return NewPolygon(points, false)
}

func NewAllowedZone() *Polygon {
	points := []Point{{60.05759504176843, 30.14495968779295},
		{60.07986463778022, 30.190278291308577},
		{60.08269008837324, 30.20143628081053},
		{60.08410272287511, 30.21662831267088},
		{60.08620015941349, 30.229760408007795},
		{60.09210650847325, 30.245724916064436},
		{60.09480253335778, 30.252677201831037},
		{60.09681370986395, 30.26031613310545},
		{60.09883394571154, 30.272926706427032},
		{60.09897148869136, 30.28416398873299},
		{60.0954811436399, 30.3286112095949},
		{60.09327428225354, 30.363445393347437},
		{60.086444857223825, 30.376478582226927},
		{60.064253238880035, 30.385181009375746},
		{60.055490095341256, 30.39468944033354},
		{60.04344323015362, 30.437052249514753},
		{60.03374429411284, 30.44212698897093},
		{60.01845570695627, 30.45914292296141},
		{60.009144281492425, 30.476695298754866},
		{59.996694566269, 30.477467774951155},
		{59.985522846219666, 30.491372346484358},
		{59.9734776331996, 30.54252743681639},
		{59.96656200178617, 30.552827119433577},
		{59.9591128611504, 30.553621053301985},
		{59.945472629965536, 30.540682077014143},
		{59.93193933304819, 30.538150071704084},
		{59.92069063807112, 30.526219606005842},
		{59.8887759685014, 30.5252754684326},
		{59.87337726427855, 30.532571076953108},
		{59.86621991030129, 30.52879452666014},
		{59.85465933529358, 30.503388642871077},
		{59.852751939209504, 30.478669404589827},
		{59.847395558876755, 30.459443330371077},
		{59.82596141511424, 30.4333508010742},
		{59.81002495078666, 30.330747589075262},
		{59.82397796859691, 30.293178893232042},
		{59.83578299691027, 30.28023103212684},
		{59.850996682094625, 30.29092323926955},
		{59.87652512736937, 30.295926020087},
		{59.88118594412356, 30.28675443686291},
		{59.88696123875598, 30.254702824938366},
		{59.89260411819843, 30.247946733331563},
		{59.89460722061685, 30.23787389382699},
		{59.90081745317471, 30.219533717360836},
		{59.903951583541954, 30.21061422347623},
		{59.906224690411726, 30.206501248146314},
		{59.90887701706451, 30.20600005944937},
		{59.92205845400411, 30.211336314284498},
		{59.93521315733946, 30.210664420926268},
		{59.946772261294086, 30.202353596293623},
		{59.966258604457494, 30.216488837802107},
		{59.976802045242714, 30.213457941615278},
		{59.98184640764717, 30.228279828631575},
		{60.00888149662998, 30.23538231810301},
		{60.02173941656657, 30.21937489470213},
		{60.03509606030931, 30.18049359282225},
		{60.04054007688507, 30.157662629687483},
		{60.049530432817626, 30.14880413604022}}

	return NewPolygon(points, true)
}
