package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	maxNodes      = 10
	volumePerPool = 255
	maxPools      = 10
	maxResources  = volumePerPool * maxPools
	maxPorts      = 16
)

func init() {
	flag.IntVar(&maxPorts, `ports`, maxPorts, `defines number of ports to generate`)
	flag.IntVar(&maxNodes, `nodes`, maxNodes, `defines number of nodes to generate`)
	flag.IntVar(&maxPools, `pools`, maxPools, `defines number of pools to generate`)
	flag.IntVar(&volumePerPool, `volumes-per-pool`, volumePerPool, `defines number of volumes per pool to generate`)
}

type (
	XmlWritable interface {
		Write(XmlWriter)
	}
)

type (
	Attribute struct {
		Name  string
		Value string
	}

	Attributes []Attribute

	InstanceAttributes struct {
		ID   string
		Rule Rule
		Attributes
	}

	MetaAttributes Attributes

	Operation struct {
		Name     string
		Interval string
		Timeout  string
	}

	Operations []Operation

	Primitive struct {
		ID                 string
		Kind               string
		MetaAttributes     MetaAttributes
		InstanceAttributes []InstanceAttributes
		Operations         Operations
	}

	Clone struct {
		Primitive      Primitive
		MetaAttributes MetaAttributes
	}

	NodeAttributeExpression struct {
		Attribute string
		Op        string
		Value     string
	}

	Expression struct {
		ID string
		NodeAttributeExpression
	}

	Expressions []Expression

	Rule struct {
		ID          string
		Score       string
		Role        string
		Op          string
		Expressions Expressions
	}

	Rules []Rule

	Location struct {
		ID         string
		ResourceID string
		Rules      Rules
	}

	Colocation struct {
		ID          string
		ResourceIDs []string
		Score       string
	}

	Action struct {
		ResourceID string
		Name       string
	}

	Actions []Action

	Order struct {
		ID      string
		Actions Actions
	}
)

type (
	UUID [5]int
)

func (uuid UUID) String() string {
	return fmt.Sprintf(`%08x-%04x-%04x-%04x-%012x`, uuid[0], uuid[1], uuid[2], uuid[3], uuid[4])
}

func (uuid UUID) ShortString() string {
	return fmt.Sprintf(`%03x-%02x`, uuid[3], uuid[4])
}

func (uuid UUID) Dup() (out UUID) {
	copy(out[:], uuid[:])
	return
}

func (c Order) Write(writer XmlWriter) {
	writer.MustStartElement(`rsc_order`)
	attributes := [][2]string{
		{`id`, c.ID},
	}
	if len(c.Actions) == 2 {
		attributes = append(attributes,
			[2]string{`first`, c.Actions[0].ResourceID},
		)
		if len(c.Actions[0].Name) > 0 {
			attributes = append(attributes,
				[2]string{`first-action`, c.Actions[0].Name},
			)
		}
		attributes = append(attributes,
			[2]string{`then`, c.Actions[1].ResourceID},
		)
		if len(c.Actions[1].Name) > 0 {
			attributes = append(attributes,
				[2]string{`then`, c.Actions[1].Name},
			)
		}
	}
	writer.MustAttributes(attributes)
	writer.MustEndElement(`rsc_order`)
}

func (c Colocation) Write(writer XmlWriter) {
	writer.MustStartElement(`rsc_colocation`)
	attributes := [][2]string{
		{`id`, c.ID},
	}
	if len(c.ResourceIDs) == 2 {
		attributes = append(attributes,
			[2]string{`rsc`, c.ResourceIDs[1]},
			[2]string{`with-rsc`, c.ResourceIDs[0]},
		)
	}
	writer.MustAttributes(attributes)
	writer.MustEndElement(`rsc_colocation`)
}

func (e NodeAttributeExpression) Write(writer XmlWriter, id string) {
	writer.MustStartElement(`expression`)
	attributes := [][2]string{
		{`id`, id},
		{`attribute`, e.Attribute},
		{`operation`, e.Op},
	}
	if len(e.Value) > 0 {
		attributes = append(attributes, [2]string{`value`, e.Value})
	}
	writer.MustAttributes(attributes)
	writer.MustEndElement(`expression`)
}

func (e Expression) Write(writer XmlWriter) {
	e.NodeAttributeExpression.Write(writer, e.ID)
}

func (e Expressions) Write(writer XmlWriter) {
	for _, e := range e {
		e.Write(writer)
	}
}

func (r Rule) Write(writer XmlWriter) {
	if len(r.Expressions) == 0 {
		return
	}
	writer.MustStartElement(`rule`)
	attributes := [][2]string{
		{`id`, r.ID},
		{`score`, r.Score},
	}
	if len(r.Op) > 0 {
		attributes = append(attributes, [2]string{`boolean-op`, r.Op})
	}
	if len(r.Role) > 0 {
		attributes = append(attributes, [2]string{`role`, r.Role})
	}
	writer.MustAttributes(attributes)
	r.Expressions.Write(writer)
	writer.MustEndElement(`rule`)
}

func (r Rules) Write(writer XmlWriter) {
	for _, r := range r {
		r.Write(writer)
	}
}

func (c Location) Write(writer XmlWriter) {
	writer.MustStartElement(`rsc_location`)
	writer.MustAttributes([][2]string{
		{`id`, c.ID},
		{`rsc`, c.ResourceID},
	})
	c.Rules.Write(writer)
	writer.MustEndElement(`rsc_location`)

}

func (attr Attribute) Write(writer XmlWriter, id string) {
	writer.MustStartElement(`nvpair`)
	writer.MustAttributes([][2]string{
		{`id`, id + `-` + attr.Name},
		{`name`, attr.Name},
		{`value`, attr.Value},
	})
	writer.MustEndElement(`nvpair`)
}

func (attrs MetaAttributes) Write(writer XmlWriter, id string) {
	if len(attrs) == 0 {
		return
	}
	writer.MustStartElement(`meta_attributes`)
	writer.MustAttributes([][2]string{
		{`id`, id + `-meta_attributes`},
	})
	for _, attr := range attrs {
		attr.Write(writer, id+`-meta_attributes`)
	}
	writer.MustEndElement(`meta_attributes`)
}

func (attrs InstanceAttributes) Write(writer XmlWriter) {
	if len(attrs.Attributes) == 0 {
		return
	}

	writer.MustStartElement(`instance_attributes`)
	writer.MustAttributes([][2]string{
		{`id`, attrs.ID},
	})
	attrs.Rule.Write(writer)
	for _, attr := range attrs.Attributes {
		attr.Write(writer, attrs.ID)
	}
	writer.MustEndElement(`instance_attributes`)
}

func (op Operation) Write(writer XmlWriter, id string) {
	writer.MustStartElement(`op`)
	writer.MustAttributes([][2]string{
		{`id`, id + `-` + op.Name + `-` + op.Interval},
		{`name`, op.Name},
		{`interval`, op.Interval},
		{`timeout`, op.Timeout},
	})
	writer.MustEndElement(`op`)
}

func (operations Operations) Write(writer XmlWriter, id string) {
	writer.MustStartElement(`operations`)
	for _, op := range operations {
		op.Write(writer, id)
	}
	writer.MustEndElement(`operations`)
}

func (p Primitive) Write(writer XmlWriter) {
	writer.MustStartElement(`primitive`)
	if len(p.Kind) == 0 {
		writer.MustAttributes([][2]string{
			{`id`, p.ID},
			{`class`, `ocf`},
			{`provider`, `heartbeat`},
			{`type`, `Dummy`},
		})
	} else {
		writer.MustAttributes([][2]string{
			{`id`, p.ID},
			{`class`, `ocf`},
			{`provider`, `yadro`},
			{`type`, p.Kind},
		})
	}

	p.MetaAttributes.Write(writer, p.ID)
	for _, attrs := range p.InstanceAttributes {
		attrs.Write(writer)
	}
	p.Operations.Write(writer, p.ID)

	writer.MustEndElement(`primitive`)
}

func (c Clone) Write(writer XmlWriter) {
	var (
		cloneID = c.Primitive.ID + `-clone`
	)

	writer.MustStartElement(`clone`)
	writer.MustAttributes([][2]string{
		{`id`, cloneID},
	})
	c.MetaAttributes.Write(writer, cloneID)
	c.Primitive.Write(writer)
	writer.MustEndElement(`clone`)
}

func main() {
	flag.Parse()

	writer := XmlWriter{Writer: os.Stdout}

	writer.MustStartElement(`cib`)
	writer.MustAttributes([][2]string{
		{`crm_feature_set`, `3.1.0`},
		{`validate-with`, `pacemaker-3.8`},
		{`epoch`, `0`},
		{`num_updates`, `0`},
		{`admin_epoch`, `0`},
		{`update-origin`, `node0`},
		{`update-client`, `cibgen`},
		{`update-user`, `root`},
		{`have-quorum`, `1`},
		{`cib-last-written`, `Thu Aug 11 16:53:06 2022`},
		{`dc-uuid`, `0`},
	})

	writer.MustStartElement(`configuration`)

	writer.MustStartElement(`crm_config`)
	writer.MustStartElement(`cluster_property_set`)
	writer.MustAttributes([][2]string{
		{`id`, `cib-bootstrap-options`},
	})
	Attribute{`stonith-enabled`, `false`}.Write(writer, `cib-bootstrap-options-stonith-enabled`)
	Attribute{`node-action-limit`, `100`}.Write(writer, `cib-bootstrap-options-node-action-limit`)
	writer.MustEndElement(`cluster_property_set`)
	writer.MustEndElement(`crm_config`)

	writeNodes(writer)

	writer.MustStartElement(`resources`)

	for poolUUID := (UUID{0, 0, 0, 0, 0}); poolUUID[3] < maxPools; poolUUID[3]++ {
		for volumeUUID := poolUUID.Dup(); volumeUUID[4] < volumePerPool; volumeUUID[4]++ {
			for portIdx := 0; portIdx < maxPorts; portIdx++ {
				mappedLun(portIdx, volumeUUID).Write(writer)
				lun(portIdx, volumeUUID).Write(writer)
			}
			volume(volumeUUID, poolUUID).Write(writer)
		}
		pool(poolUUID).Write(writer)
	}

	for aclIdx := 0; aclIdx < maxPorts; aclIdx++ {
		acl(aclIdx).Write(writer)
	}

	for targetIdx := 0; targetIdx < maxPorts; targetIdx++ {
		target(targetIdx).Write(writer)
	}

	for portIdx := 0; portIdx < maxPorts; portIdx++ {
		port(portIdx).Write(writer)
	}

	traidConfig().Write(writer)

	writer.MustEndElement(`resources`)

	writer.MustStartElement(`constraints`)

	for poolUUID := (UUID{0, 0, 0, 0, 0}); poolUUID[3] < maxPools; poolUUID[3]++ {
		for volumeUUID := poolUUID.Dup(); volumeUUID[4] < volumePerPool; volumeUUID[4]++ {

			for portIdx := 0; portIdx < maxPorts; portIdx++ {
				for _, constr := range mappedLunConstraints(portIdx, volumeUUID) {
					constr.Write(writer)
				}

				for _, constr := range lunConstraints(portIdx, volumeUUID) {
					constr.Write(writer)
				}
			}

			for _, constr := range volumeConstraints(volumeUUID, poolUUID) {
				constr.Write(writer)
			}
		}
		for _, constr := range poolConstraints(poolUUID) {
			constr.Write(writer)
		}
	}

	for aclIdx := 0; aclIdx < maxPorts; aclIdx++ {
		for _, constr := range aclConstraints(aclIdx) {
			constr.Write(writer)
		}
	}

	for targetIdx := 0; targetIdx < maxPorts; targetIdx++ {
		for _, constr := range targetConstraints(targetIdx) {
			constr.Write(writer)
		}
	}

	for portIdx := 0; portIdx < maxPorts; portIdx++ {
		for _, constr := range portConstraints(portIdx) {
			constr.Write(writer)
		}
	}

	for _, constr := range traidConfigConstraints() {
		constr.Write(writer)
	}
	writer.MustEndElement(`constraints`)

	writer.MustEndElement(`configuration`)

	writer.MustEndElement(`cib`)
}

func portIDFromIdx(idx int) string {
	return fmt.Sprintf(`p%02x`, idx)
}

func targetIDFromIdx(idx int) string {
	return fmt.Sprintf(`tgt-%02x`, idx)
}

func aclIDFromIdx(idx int) string {
	return fmt.Sprintf(`acl-%02x`, idx)
}

func lunIDFromIdxAndVolumeUUID(idx int, uuid UUID) string {
	return fmt.Sprintf(`lun-%s-%02x`, uuid.ShortString(), idx)
}

func mappedLunIDFromIdxAndVolumeUUID(idx int, uuid UUID) string {
	return fmt.Sprintf(`mlun-%s-%02x`, uuid.ShortString(), idx)
}

func mappedLunConstraints(idx int, uuid UUID) []XmlWritable {
	var (
		aclID       = aclIDFromIdx(idx)
		lunID       = lunIDFromIdxAndVolumeUUID(idx, uuid)
		mappedLunID = mappedLunIDFromIdxAndVolumeUUID(idx, uuid)
	)
	return []XmlWritable{
		Colocation{
			ID:          mappedLunID + `-with-` + aclID,
			ResourceIDs: []string{aclID + `-clone`, mappedLunID + `-clone`},
			Score:       `INFINITY`,
		},
		Colocation{
			ID:          mappedLunID + `-with-` + lunID,
			ResourceIDs: []string{lunID + `-clone`, mappedLunID + `-clone`},
			Score:       `INFINITY`,
		},
		Order{
			ID: mappedLunID + `-after-` + aclID,
			Actions: Actions{{
				ResourceID: aclID + `-clone`,
			}, {
				ResourceID: mappedLunID + `-clone`,
			}},
		},
		Order{
			ID: mappedLunID + `-after-` + lunID,
			Actions: Actions{{
				ResourceID: lunID + `-clone`,
			}, {
				ResourceID: mappedLunID + `-clone`,
			}},
		},
	}
}

func mappedLun(idx int, uuid UUID) XmlWritable {
	var (
		mappedLunID = mappedLunIDFromIdxAndVolumeUUID(idx, uuid)
	)
	return Clone{
		Primitive{
			mappedLunID,
			`mapped-lun`,
			MetaAttributes{},
			[]InstanceAttributes{{
				ID: mappedLunID + `-instance_attributes`,
				Attributes: Attributes{{
					`lun`, lunIDFromIdxAndVolumeUUID(idx, uuid),
				}},
			}},
			Operations{
				{`stop`, `0s`, `15s`},
				{`start`, `0s`, `60s`},
				{`monitor`, `30s`, `10s`},
			},
		},
		MetaAttributes{
			{`target-role`, `Started`},
			{`interleave`, `true`},
		},
	}
}

func lunConstraints(idx int, uuid UUID) []XmlWritable {
	var (
		targetID = targetIDFromIdx(idx)
		lunID    = lunIDFromIdxAndVolumeUUID(idx, uuid)
		volumeID = volumeIDFromUUID(uuid)
	)
	return []XmlWritable{
		Colocation{
			ID:          lunID + `-with-` + targetID,
			ResourceIDs: []string{targetID + `-clone`, lunID + `-clone`},
			Score:       `INFINITY`,
		},
		Colocation{
			ID:          lunID + `-with-` + volumeID,
			ResourceIDs: []string{volumeID + `-clone`, lunID + `-clone`},
			Score:       `INFINITY`,
		},
		Order{
			ID: lunID + `-after-` + targetID,
			Actions: Actions{{
				ResourceID: targetID + `-clone`,
			}, {
				ResourceID: lunID + `-clone`,
			}},
		},
		Order{
			ID: lunID + `-after-` + volumeID,
			Actions: Actions{{
				ResourceID: volumeID + `-clone`,
			}, {
				ResourceID: lunID + `-clone`,
			}},
		},
	}
}

func lun(idx int, uuid UUID) XmlWritable {
	var (
		lunID = lunIDFromIdxAndVolumeUUID(idx, uuid)
	)
	return Clone{
		Primitive{
			lunID,
			`lun`,
			MetaAttributes{},
			[]InstanceAttributes{{
				ID: lunID + `-instance_attributes`,
				Attributes: Attributes{{
					`volume`, volumeIDFromUUID(uuid),
				}},
			}},
			Operations{
				{`stop`, `0s`, `15s`},
				{`start`, `0s`, `60s`},
				{`monitor`, `30s`, `10s`},
			},
		},
		MetaAttributes{
			{`target-role`, `Started`},
			{`interleave`, `true`},
		},
	}
}

func aclConstraints(idx int) []XmlWritable {
	var (
		aclID    = aclIDFromIdx(idx)
		targetID = targetIDFromIdx(idx)
	)
	return []XmlWritable{
		Colocation{
			ID:          aclID + `-with-` + targetID,
			ResourceIDs: []string{targetID + `-clone`, aclID + `-clone`},
			Score:       `INFINITY`,
		},
		Order{
			ID: aclID + `-after-` + targetID,
			Actions: Actions{{
				ResourceID: targetID + `-clone`,
			}, {
				ResourceID: aclID + `-clone`,
			}},
		},
	}
}

func acl(idx int) XmlWritable {
	var (
		aclID = aclIDFromIdx(idx)
	)
	return Clone{
		Primitive{
			aclID,
			`acl`,
			MetaAttributes{},
			[]InstanceAttributes{{
				ID: aclID + `-instance_attributes`,
				Attributes: Attributes{{
					`port`, portIDFromIdx(idx),
				}},
			}},
			Operations{
				{`stop`, `0s`, `15s`},
				{`start`, `0s`, `60s`},
				{`monitor`, `30s`, `10s`},
			},
		},
		MetaAttributes{
			{`target-role`, `Started`},
			{`interleave`, `true`},
		},
	}
}

func targetConstraints(idx int) []XmlWritable {
	var (
		targetID = targetIDFromIdx(idx)
		portID   = portIDFromIdx(idx)
	)
	return []XmlWritable{
		Colocation{
			ID:          targetID + `-with-` + portID,
			ResourceIDs: []string{portID + `-clone`, targetID + `-clone`},
			Score:       `INFINITY`,
		},
		Order{
			ID: targetID + `-after-` + portID,
			Actions: Actions{{
				ResourceID: portID + `-clone`,
			}, {
				ResourceID: targetID + `-clone`,
			}},
		},
	}
}

func target(idx int) XmlWritable {
	var (
		targetID = targetIDFromIdx(idx)
	)
	return Clone{
		Primitive{
			targetID,
			`target`,
			MetaAttributes{},
			[]InstanceAttributes{{
				ID: targetID + `-instance_attributes`,
				Attributes: Attributes{{
					`port`, portIDFromIdx(idx),
				}},
			}},
			Operations{
				{`stop`, `0s`, `15s`},
				{`start`, `0s`, `60s`},
				{`monitor`, `30s`, `10s`},
			},
		},
		MetaAttributes{
			{`target-role`, `Started`},
			{`interleave`, `true`},
		},
	}
}

func portConstraints(idx int) []XmlWritable {
	return []XmlWritable{
		locationOnlyOnStorageProcessor(portIDFromIdx(idx)),
	}
}

func port(idx int) XmlWritable {
	var (
		portID             = portIDFromIdx(idx)
		instanceAttributes = make([]InstanceAttributes, maxNodes)
	)

	for nodeIdx := 0; nodeIdx < maxNodes; nodeIdx++ {
		id := portID + `-instance_attributes-` + strconv.Itoa(nodeIdx)
		instanceAttributes[nodeIdx] = InstanceAttributes{
			ID: id,
			Attributes: Attributes{
				{`address`, fmt.Sprintf(`192.168.%d.1%02d/16`, nodeIdx, idx)},
			},
			Rule: Rule{
				ID:    id + `-rule-0`,
				Score: `0`,
				Expressions: Expressions{{
					ID: id + `-rule-0-expression-0`,
					NodeAttributeExpression: NodeAttributeExpression{
						Attribute: `#uname`,
						Value:     `node` + strconv.Itoa(nodeIdx),
						Op:        `eq`,
					},
				}},
			},
		}
	}

	return Clone{
		Primitive{
			portID,
			`port`,
			MetaAttributes{},
			instanceAttributes,
			Operations{
				{`stop`, `0s`, `15s`},
				{`start`, `0s`, `60s`},
				{`monitor`, `30s`, `10s`},
			},
		},
		MetaAttributes{
			{`target-role`, `Started`},
			{`interleave`, `true`},
		},
	}
}

func locationOnlyOnStorageProcessor(id string) XmlWritable {
	return Location{
		ID:         id + `-only-on-storage-processors`,
		ResourceID: id + `-clone`,
		Rules: Rules{{
			ID:    id + `-only-on-storage-processors-rule-0`,
			Score: `-INFINITY`,
			Op:    `or`,
			Expressions: Expressions{
				{
					ID: id + `-only-on-storage-processors-rule-0-expression-0`,
					NodeAttributeExpression: NodeAttributeExpression{
						Attribute: `node-type`,
						Op:        `not_defined`,
					},
				},
				{
					ID: id + `-only-on-storage-processors-rule-0-expression-1`,
					NodeAttributeExpression: NodeAttributeExpression{
						Attribute: `node-type`,
						Op:        `ne`,
						Value:     `storage-processor`,
					},
				},
			},
		}},
	}
}

func volumeConstraints(volumeUUID, poolUUID UUID) []XmlWritable {
	var (
		volumeID = volumeIDFromUUID(volumeUUID)
		poolID   = poolIDFromUUID(poolUUID)
	)
	return []XmlWritable{
		Colocation{
			ID:          volumeID + `-with-` + poolID,
			ResourceIDs: []string{poolID + `-clone`, volumeID + `-clone`},
			Score:       `INFINITY`,
		},
		Order{
			ID: volumeID + `-after-` + poolID,
			Actions: Actions{{
				ResourceID: poolID + `-clone`,
			}, {
				ResourceID: volumeID + `-clone`,
			}},
		},
	}
}

func volumeIDFromUUID(uuid UUID) string {
	return fmt.Sprintf(`vol-%s`, uuid.ShortString())
}

func volume(volumeUUID, poolUUID UUID) XmlWritable {
	var (
		volumeID = volumeIDFromUUID(volumeUUID)
	)

	return Clone{
		Primitive{
			volumeID,
			`volume`,
			MetaAttributes{},
			[]InstanceAttributes{{
				ID: volumeID + `-instance_attributes`,
				Attributes: Attributes{
					{`name`, volumeID},
					{`uuid`, volumeUUID.String()},
					{`pool`, `pool-` + poolUUID.String()},
				},
			}},
			Operations{
				{`stop`, `0s`, `15s`},
				{`start`, `0s`, `60s`},
				{`monitor`, `30s`, `10s`},
			},
		},
		MetaAttributes{
			{`target-role`, `Started`},
			{`interleave`, `true`},
		},
	}
}

func poolIDFromUUID(uuid UUID) string {
	return fmt.Sprintf(`pool-%s`, uuid.ShortString())
}

func poolConstraints(uuid UUID) []XmlWritable {
	var (
		poolID = poolIDFromUUID(uuid)
	)
	return []XmlWritable{
		Colocation{
			ID:          poolID + `-with-traid-config`,
			ResourceIDs: []string{`traid-config-clone`, poolID + `-clone`},
			Score:       `INFINITY`,
		},
		Order{
			ID: poolID + `-after-traid-config`,
			Actions: Actions{{
				ResourceID: `traid-config-clone`,
			}, {
				ResourceID: poolID + `-clone`,
			}},
		},
	}
}

func pool(uuid UUID) XmlWritable {
	var (
		uuidStr = uuid.String()
		poolID  = poolIDFromUUID(uuid)
	)
	return Clone{
		Primitive{
			poolID,
			`pool`,
			MetaAttributes{},
			[]InstanceAttributes{{
				ID: poolID + `-instance_attributes`,
				Attributes: Attributes{
					{`name`, poolID},
					{`uuid`, uuidStr},
				},
			}},
			Operations{
				{`stop`, `0s`, `15s`},
				{`start`, `0s`, `60s`},
				{`monitor`, `30s`, `10s`},
			},
		},
		MetaAttributes{
			{`target-role`, `Started`},
			{`interleave`, `true`},
		},
	}
}

func traidConfigConstraints() []XmlWritable {
	return []XmlWritable{
		Location{
			ID:         `traid-config-only-on-storage-processors`,
			ResourceID: `traid-config-clone`,
			Rules: Rules{{
				ID:    `traid-config-only-on-storage-processors-rule-0`,
				Score: `INFINITY`,
				Role:  `Master`,
				Expressions: Expressions{
					{
						ID: `traid-config-only-on-storage-processors-rule-0-expression-0`,
						NodeAttributeExpression: NodeAttributeExpression{
							Attribute: `node-type`,
							Op:        `eq`,
							Value:     `storage-processor`,
						},
					},
				},
			}, {
				ID:    `traid-config-only-on-storage-processors-rule-1`,
				Score: `-INFINITY`,
				Op:    `or`,
				Expressions: Expressions{
					{
						ID: `traid-config-only-on-storage-processors-rule-1-expression-0`,
						NodeAttributeExpression: NodeAttributeExpression{
							Attribute: `node-type`,
							Op:        `not_defined`,
						},
					},
					{
						ID: `traid-config-only-on-storage-processors-rule-1-expression-1`,
						NodeAttributeExpression: NodeAttributeExpression{
							Attribute: `node-type`,
							Op:        `ne`,
							Value:     `storage-processor`,
						},
					},
				},
			}},
		},
	}
}
func traidConfig() XmlWritable {
	return Clone{
		Primitive{
			`traid-config`,
			`traid-config`,
			MetaAttributes{},
			[]InstanceAttributes{{
				ID: `traid-config-instance_attributes`,
				Attributes: Attributes{
					{`active`, `true`},
				},
			}},
			Operations{
				{`stop`, `0s`, `60s`},
				{`start`, `0s`, `60s`},
				{`reload`, `0s`, `60s`},
				{`monitor`, `20s`, `20s`},
				{`monitor`, `10s`, `20s`},
				{`notify`, `0s`, `30s`},
				{`promote`, `0s`, `30s`},
				{`demote`, `0s`, `30s`},
			},
		},
		MetaAttributes{
			{`promotable`, `true`},
			{`notify`, `true`},
			{`target-role`, `Started`},
			{`priority`, `999999`},
			{`failure-timeout`, `5m`},
		},
	}
}

func writeNodes(writer XmlWriter) {
	writer.MustStartElement(`nodes`)
	for nodeIndex := 0; nodeIndex < maxNodes; nodeIndex++ {
		writeNode(writer, nodeIndex)
	}
	writer.MustEndElement(`nodes`)
}

func writeNode(writer XmlWriter, nodeIndex int) {
	writer.MustStartElement(`node`)
	writer.MustAttributes([][2]string{
		{`id`, strconv.Itoa(nodeIndex)},
		{`uname`, `node` + strconv.Itoa(nodeIndex)},
	})
	InstanceAttributes{
		ID:         `nodes-` + strconv.Itoa(nodeIndex),
		Attributes: Attributes{{`node-type`, `storage-processor`}},
	}.Write(writer)
	writer.MustEndElement(`node`)
}
