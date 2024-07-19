package player

import (
	"sync"

	"github.com/dynamitemc/aether/atomic"
	"github.com/dynamitemc/aether/net/metadata"
	"github.com/dynamitemc/aether/net/packet/configuration"
	"github.com/dynamitemc/aether/server/entity"
)

var _ entity.Entity = (*Player)(nil)

type Player struct {
	entityId int32

	x, y, z    atomic.AtomicValue[float64]
	yaw, pitch atomic.AtomicValue[float32]

	clientInfo atomic.AtomicValue[configuration.ClientInformation]

	md_mu    sync.Mutex
	metadata metadata.Metadata
}

// NewPlayer creates a player struct with the entity id specified and initalizes an entity metadata map for it
func NewPlayer(entityId int32) *Player {
	return &Player{
		entityId: entityId,
		metadata: metadata.Metadata{
			// Entity
			metadata.BaseIndex:                      metadata.Byte(0),
			metadata.AirTicksIndex:                  metadata.VarInt(300),
			metadata.CustomNameIndex:                metadata.OptionalTextComponent(nil),
			metadata.IsCustomNameVisibleIndex:       metadata.Boolean(false),
			metadata.IsSilentIndex:                  metadata.Boolean(false),
			metadata.HasNoGravityIndex:              metadata.Boolean(false),
			metadata.PoseIndex:                      metadata.Standing,
			metadata.TicksFrozenInPowderedSnowIndex: metadata.VarInt(0),
			// Player extends Living Entity
			metadata.PlayerAdditionalHeartsIndex:   metadata.Float(0),
			metadata.PlayerScoreIndex:              metadata.VarInt(0),
			metadata.PlayerDisplayedSkinPartsIndex: metadata.Byte(0),
			metadata.PlayerMainHandIndex:           metadata.Byte(1),
		},
	}
}

func (p *Player) Position() (x, y, z float64) {
	return p.x.Get(), p.y.Get(), p.z.Get()
}

func (p *Player) Rotation() (yaw, pitch float32) {
	return p.yaw.Get(), p.pitch.Get()
}

func (p *Player) SetPosition(x, y, z float64) {
	p.x.Set(x)
	p.y.Set(y)
	p.z.Set(z)
}

func (p *Player) SetRotation(yaw, pitch float32) {
	p.yaw.Set(yaw)
	p.pitch.Set(pitch)
}

func (p *Player) EntityId() int32 {
	return p.entityId
}

func (p *Player) SetClientInformation(info configuration.ClientInformation) {
	p.clientInfo.Set(info)
}

func (p *Player) ClientInformation() configuration.ClientInformation {
	return p.clientInfo.Get()
}

func (p *Player) Metadata() metadata.Metadata {
	p.md_mu.Lock()
	defer p.md_mu.Unlock()
	return p.metadata
}

func (p *Player) SetMetadata(md metadata.Metadata) {
	p.md_mu.Lock()
	defer p.md_mu.Unlock()
	p.metadata = md
}

func (p *Player) MetadataIndex(i byte) any {
	p.md_mu.Lock()
	defer p.md_mu.Unlock()
	return p.metadata[i]
}

func (p *Player) SetMetadataIndex(i byte, v any) {
	p.md_mu.Lock()
	defer p.md_mu.Unlock()
	p.metadata[i] = v
}

func (p *Player) SetMetadataIndexes(md metadata.Metadata) {
	p.md_mu.Lock()
	defer p.md_mu.Unlock()
	for index, value := range md {
		p.metadata[index] = value
	}
}
