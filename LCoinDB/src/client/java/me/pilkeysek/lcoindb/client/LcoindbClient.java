package me.pilkeysek.lcoindb.client;

import me.pilkeysek.lcoindb.client.command.LCoinCommand;
import net.fabricmc.api.ClientModInitializer;
import net.fabricmc.fabric.api.client.command.v2.ClientCommandRegistrationCallback;
import net.fabricmc.fabric.api.client.networking.v1.ClientPlayConnectionEvents;
import net.minecraft.text.Text;
import net.minecraft.util.Formatting;

import java.util.List;
import java.util.Objects;
import java.util.UUID;

public class LcoindbClient implements ClientModInitializer {
    public static final me.pilkeysek.lcoindb.client.MainConfig config = me.pilkeysek.lcoindb.client.MainConfig.createAndLoad();

    @Override
    public void onInitializeClient() {
        if(config.secret().isEmpty() || config.secret() == null) {
            config.secret(UUID.randomUUID().toString().replace("-","")); // The authentication key is a random uuid, yes
        }
        ClientPlayConnectionEvents.JOIN.register( (clientPlayNetworkHandler, packetSender, minecraftClient) -> {
            if(!config.authenticationEnabled()) return;
            if(!minecraftClient.isConnectedToLocalServer()) {
                assert minecraftClient.player != null;
                String serverIP = Objects.requireNonNull(minecraftClient.player.networkHandler.getServerInfo()).address;
                System.out.println("SERVER IP: " + serverIP);
                List<String> trustedAuthServers = config.trustedAuthServers();
                if(trustedAuthServers.contains(serverIP)) {
                    minecraftClient.player.sendMessage(Text.literal("Attempting to authenticate ...\nYou may have to rejoin a second time if it doesn't succeed immediately (it might give an error)").formatted(Formatting.AQUA));
                    minecraftClient.player.networkHandler.sendChatCommand("authme " + config.secret());
                }
            }
        });
        ClientCommandRegistrationCallback.EVENT.register(((commandDispatcher, commandRegistryAccess) -> {
            LCoinCommand.register(commandDispatcher);
        }));
    }
}